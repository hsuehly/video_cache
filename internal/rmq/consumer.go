package rmq

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
	"video_cache/domain"
	"video_cache/internal/m3task"
	"video_cache/internal/redisscript"
	"video_cache/pkg/feishu"
	"video_cache/pkg/httpclient"
	"video_cache/pkg/logger"

	"video_cache/bootstrap"
)

type UPDate struct {
	Pan      int    `db:"pan"`
	PathId   string `db:"path_id"`
	FileId   string `db:"file_id"`
	LinkMD5  string `db:"link_md5"`
	Link     string `db:"link"`
	FileName string `db:"file_name"`
}

// XGROUP CREATE m3u8_topic m3u8_group 0-0
// XGROUP CREATE m3u8_topic m3u8_group $
// 接收到消息后执行的回调函数
//type MsgCallback func(ctx context.Context, msg *MsgEntity) error

// 消费者
type Consumer struct {
	// consumer 生命周期管理
	ctx  context.Context
	stop context.CancelFunc

	// 接收到 msg 时执行的回调函数，由使用方定义
	//callbackFunc MsgCallback

	// redis 客户端，基于 redis 实现 message queue
	client *redis.Client

	// 消费的 topic
	topic string
	// 所属的消费者组
	groupID string
	// 当前节点的消费者 id
	consumerID string

	// 各消息累计失败次数
	//failureCnts map[MsgEntity]int
	httpClient *httpclient.HttpClient
	parse      *m3task.Parse
	app        *bootstrap.Application
	// 一些用户自定义的配置
	opts *ConsumerOptions
}

func NewConsumer(app *bootstrap.Application, opts ...ConsumerOption) (*Consumer, error) {

	ctx, stop := context.WithCancel(context.Background())
	c := Consumer{
		client:     app.Rdb,
		ctx:        ctx,
		stop:       stop,
		topic:      app.Env.ToPic,
		groupID:    app.Env.ConsumerGroup,
		consumerID: app.Env.ConsumerID,
		httpClient: httpclient.NewHttpClient(),
		parse:      m3task.NewParse(),
		app:        app,
		opts:       &ConsumerOptions{},

		//failureCnts: make(map[MsgEntity]int),
	}

	if err := c.checkParam(); err != nil {
		return nil, err
	}

	for _, opt := range opts {
		opt(c.opts)
	}

	repairConsumer(c.opts)

	go c.run()
	return &c, nil
}

func (c *Consumer) checkParam() error {
	//if c.callbackFunc == nil {
	//	return errors.New("callback function can't be empty")
	//}

	if c.client == nil {
		return errors.New("redis client can't be empty")
	}

	if c.topic == "" || c.consumerID == "" || c.groupID == "" {
		return errors.New("topic | group_id | consumer_id can't be empty")
	}

	return nil
}

// 停止 consumer

func (c *Consumer) Stop() {
	c.stop()
}

// 运行消费者
func (c *Consumer) run() {
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
		}
		// 新消息接收处理
		msgs, err := c.receive()
		if err != nil {
			continue
		}
		//if err != nil && !errors.Is(err, redis.Nil) {
		//	continue
		//}
		c.handlerMsgs(msgs)

		//if msgs != nil {
		//	c.handlerMsgs(msgs)
		//}

		pendingMsgs, err := c.receivePending()
		if err != nil {
			continue
		}
		//if err != nil && !errors.Is(err, redis.Nil) {
		//	continue
		//}
		c.handlerMsgs(pendingMsgs)
		//if pendingMsgs != nil {
		//	c.handlerMsgs(pendingMsgs)
		//}
	}
}

func (c *Consumer) receive() (*redis.XMessage, error) {
	return c.xReadGroup(int(c.opts.receiveTimeout.Milliseconds()), false)
}

func (c *Consumer) receivePending() (*redis.XMessage, error) {
	return c.xReadGroup(0, true)
}
func (c *Consumer) handlerMsgs(messages *redis.XMessage) {
	defer func(id string) {
		if err := c.xACK(id); err != nil {
			logger.Warn("msg ack failed", logger.String("msg id:", messages.ID), logger.Err(err))
		}
	}(messages.ID)
	//defer func() {
	//if messages != nil {
	//	if err := c.xACK(ctx, messages.ID); err != nil {
	//		logger.Warn("msg ack failed", logger.String("msg id:", messages.ID), logger.Err(err))
	//	}
	//ids := make([]string, 0, 10)
	v, ok := messages.Values["task"]
	if !ok {
		logger.Warn("读取消息失败")
		return
	}
	var Data = new(domain.RmqPostData)
	value, ok := v.(string)
	if !ok {
		logger.Warn("断言消息失败")
		return
	}
	// callback 执行成功，进行 ack
	err := json.Unmarshal([]byte(value), Data)
	if err != nil {
		logger.Warn("序列化失败", logger.Err(err))
		return
	}
	lastDotIndex := strings.LastIndex(Data.FileName, ".")
	if lastDotIndex == -1 {
		logger.Warn("解析redKey错误：", logger.String("资源", Data.Url), logger.Err(err))
		return
	}
	redKey := Data.FileName[:lastDotIndex]
	defer func() {
		reply, err := c.HDel(redKey, "task")
		if err != nil {
			logger.Warn("清除任务列表错误", logger.Err(err))
		}
		if reply == 0 {
			logger.Warn("清除任务列表task字段失败返回值为0")
		}
	}()
	resp, err := c.httpClient.Do("GET", Data.Surl, nil, true, nil)
	if err != nil {
		logger.Warn("请求错误", logger.Err(err))
		return
	}
	if resp.StatusCode != 200 {
		logger.Warn("资源错误")
		return
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(resp.Body)
	resp.Body.Close()
	var m3u8Data []byte
	switch Data.TaskId {
	case 0:
		m3u8Data, err = c.parse.M3(buf, Data.BaseUrl)
	case 1:
		m3u8Data, err = c.parse.ZA(buf, Data.BaseUrl)
	case 2:
		m3u8Data, err = c.parse.AiKu(buf, Data.BaseUrl)
	}
	if err != nil {
		//rely, _ := c.HIncrBy(redKey, "ecount", 1)
		_ = feishu.SendMsg(fmt.Sprintf("解析m3u8 失败url: %s surl:%s 错误原因: %s", Data.Url, Data.Surl, err.Error()))
		logger.Warn("解析m3u8失败：", logger.String("url", Data.Url), logger.String("surl", Data.Surl), logger.Err(err))
		return
	}
	fileId, err := c.app.YunPan139.Put(Data.PathIds, Data.FileName, m3u8Data)
	if err != nil {
		logger.Warn("上传云盘出错：", logger.String("资源", Data.Url), logger.Err(err))
		_ = feishu.SendMsg(fmt.Sprintf("上传云盘出错资源：%s,错误详情%s", Data.Url, err))
		return
	}
	data := &UPDate{
		Pan:      3,
		LinkMD5:  Data.LinkMd5,
		Link:     Data.Url,
		FileId:   fileId,
		PathId:   Data.PathIds,
		FileName: Data.FileName,
	}
	err = c.saveM3(Data.TableName, data)
	if err != nil {
		logger.Warn("保存mysql出错：", logger.String("资源", Data.Url), logger.Err(err))
		_ = feishu.SendMsg(fmt.Sprintf("保存mysql出错：%s,错误详情%s", Data.Url, err))
		err = c.app.YunPan139.Remove(fileId)
		if err != nil {
			_ = feishu.SendMsg(fmt.Sprintf("mysql出错删除云盘资源错误：资源id %s,资源链接 %s 错误详情%s", fileId, Data.Url, err))
		}
		return
	}
	//if err = c.Del(redKey); err != nil {
	//	logger.Warn("删除缓存失败：", logger.String("资源", Data.Url), logger.Err(err))
	//}
	return
}
func (c *Consumer) xACK(msgID ...string) error {
	if msgID == nil {
		return errors.New("redis xACK msg_ id can't be empty")
	}
	reply, err := c.client.XAck(c.ctx, c.topic, c.groupID, msgID...).Result()
	if err != nil {
		return err
	}
	if reply == 0 {
		return fmt.Errorf("invalid reply: %d", reply)
	}
	return nil
}
func (c *Consumer) Del(key string) error {
	err := c.client.Del(c.ctx, key).Err()
	//reply := c.client.XAck(ctx, c.topic, c.groupID, msgID).Val()
	if err != nil {
		return err
	}

	return nil
}
func (c *Consumer) Incr(key string) (int64, error) {
	reply, err := redisscript.IncrTimeScript.Run(c.ctx, c.client, []string{key}, 7200).Result()
	//reply := c.client.XAck(ctx, c.topic, c.groupID, msgID).Val()
	if err != nil {
		return 0, err
	}
	num, ok := reply.(int64)
	if !ok {
		return 0, errors.New("错误队列redis断言失败")
	}
	return num, nil
}

// HSetTimeScript
func (c *Consumer) HSet(key, url string) error {
	return redisscript.HSetTimeScript.Run(c.ctx, c.client, []string{key}, url, "m3u8", "1", 3600).Err()
}

func (c *Consumer) xReadGroup(timeoutMiliSeconds int, pending bool) (*redis.XMessage, error) {
	//conn, err := c.pool.GetContext(ctx)
	//if err != nil {
	//	return nil, err
	//}
	//defer conn.Close()

	// consumer 刚启动时，批量获取一次分配给本节点，但是还没 ack 的消息进行处理
	// consumer 处理消息之后，如果想给一个坏的 ack，那则是再获取一次 pending 重新走一次流程
	// 分配给本节点，但是尚未 ack 的消息 0-0
	// 拿到尚未分配过的新消息 >
	// [{1713712293762-0 map[]} {1713790835904-0 map[]} {1713793879331-0 map[hsuehly:454545]}]}]
	var rawReply []redis.XStream
	var err error
	//resu := c.client.HGetAll(ctx, "462e3f24d1a6be473bbdb045b81a7885")
	//fmt.Println(resu, "resu")
	var opt = &redis.XReadGroupArgs{
		Group:    c.groupID,
		Consumer: c.consumerID,
		Count:    1,
	}
	if pending {
		opt.Streams = []string{c.topic, "0-0"}
		rawReply, err = c.client.XReadGroup(c.ctx, opt).Result()
	} else {
		opt.Streams = []string{c.topic, ">"}
		opt.Block = time.Duration(timeoutMiliSeconds) * time.Millisecond
		rawReply, err = c.client.XReadGroup(c.ctx, opt).Result()
	}
	if err != nil {
		return nil, err
	}
	if len(rawReply) == 0 {
		return nil, redis.Nil
	}
	replyElement := rawReply[0].Messages
	if len(replyElement) == 0 {
		return nil, redis.Nil
	}
	//return &replyElement[0], nil
	//fmt.Println(replyElement[0], "replyElement[0]")
	return &replyElement[0], nil
}

// saveMysql
func (c *Consumer) saveM3(tableName string, data *UPDate) error {
	ADM3U8SQL := "INSERT INTO " + tableName + " (pan, link, link_md5, path_id, file_id ,file_name) VALUES (:pan, :link, :link_md5, :path_id, :file_id , :file_name) ON DUPLICATE KEY UPDATE pan = :pan, path_id = :path_id, file_id = :file_id, file_name = :file_name, status = 0"
	_, err := c.app.DB.NamedExecContext(c.ctx, ADM3U8SQL, data)
	if err != nil {
		return err
	}
	return nil
}

//	func (c *Consumer) errHandle(urlKey, redKey, surl string) {
//		errKey := "task_err:" + urlKey
//		num, err := c.Incr(errKey)
//		if err != nil {
//			logger.Warn("添加错误key出错：", logger.Err(err))
//		}
//		if num >= 3 {
//			if err = c.Del(errKey); err != nil {
//				logger.Warn("删除错误key出错：", logger.Err(err))
//			}
//			if err = c.HSet(redKey, surl); err != nil {
//				logger.Warn("错误3次设置缓存错误：", logger.Err(err))
//			}
//			_ = feishu.SendMsg(fmt.Sprintf("资源重试到达三次 %s ", surl))
//		}
//	}
func (c *Consumer) HDel(key string, field ...string) (int64, error) {
	return c.client.HDel(c.ctx, key, field...).Result()
}
func (c *Consumer) HIncrBy(key string, field string, incr int64) (int64, error) {
	return c.client.HIncrBy(c.ctx, key, field, incr).Result()
}
