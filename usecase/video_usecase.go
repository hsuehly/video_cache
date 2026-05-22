package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"video_cache/bootstrap"
	"video_cache/domain"
	"video_cache/internal/fetcherr"
	"video_cache/internal/rmq"
	"video_cache/internal/urlutil"
	"video_cache/pkg/feishu"
	"video_cache/pkg/logger"
	"video_cache/pkg/yun139"

	"github.com/goccy/go-json"
)

type videoUsecase struct {
	videoRepository domain.VideoRepository
	env             *bootstrap.Env
	producer        *rmq.Producer
	contextTimeout  time.Duration
	yunPan139       *yun139.Yun139
}

func NewVideoUsecase(videoRepository domain.VideoRepository, env *bootstrap.Env, producer *rmq.Producer, yunPan139 *yun139.Yun139, timeout time.Duration) domain.VideoUsecase {
	return &videoUsecase{
		videoRepository: videoRepository,
		env:             env,
		producer:        producer,
		yunPan139:       yunPan139,
		contextTimeout:  timeout,
	}
}
func (vu *videoUsecase) GetLink(c context.Context, url string) (*domain.GetLinkResp, error) {
	ctx, cancel := context.WithTimeout(c, vu.contextTimeout)
	defer cancel()
	key := urlutil.EncodeMD5WithSalt(url, []byte(vu.env.CacheKey))
	Url, Type, TTl := vu.videoRepository.HGetUrlTypeTime(ctx, key)
	if Url == "" {
		linkMd5 := urlutil.EncodeMD5(url)
		table := urlutil.GetTable(url)
		if table == "" {
			return nil, fetcherr.TableNoLink
		}
		link, err := vu.videoRepository.GetLinkByMd5(ctx, table, linkMd5)
		if err != nil {
			return nil, err
		}
		switch link.Pan {
		case 3:
			Url, err = vu.yunPan139.GetLink(link.FileId)
			TTl = 25200
			Type = "m3u8"
		}
		if err != nil {
			if errors.Is(err, yun139.NoFileErr) {
				sqlErr := vu.videoRepository.DelByLinkMd5(ctx, table, linkMd5)
				if sqlErr != nil {
					return nil, sqlErr
				}
			}
			return nil, err
		}
		if Url == "" {
			return nil, errors.New("获取链接为空")
		}
		err = vu.videoRepository.HSetTimeScript(ctx, key, Url, Type, "0", TTl)
		if err != nil {
			logger.Warn("url添加缓存失败", logger.Err(err))
		}
		return &domain.GetLinkResp{Url: Url, Type: Type, ExpTime: TTl}, nil
	}
	return &domain.GetLinkResp{Url: Url, Type: Type, ExpTime: TTl}, nil
}

func (vu *videoUsecase) TaskM3(c context.Context, body *domain.TaskM3Body) (string, error) {
	ctx, cancel := context.WithTimeout(c, vu.contextTimeout)
	defer cancel()
	key := urlutil.EncodeMD5WithSalt(body.Url, []byte(vu.env.CacheKey))

	if body.Time != 0 {
		if err := vu.videoRepository.HSetTimeScript(ctx, key, body.Surl, "m3u8", "1", body.Time); err != nil {
			return "", err
		}
		return "cache ok", nil
	}
	if err := vu.videoRepository.HSetTimeScript(ctx, key, body.Surl, "m3u8", "1", 3600); err != nil {
		return "", err
	}
	linkMd5 := urlutil.EncodeMD5(body.Url)
	//taskKey := "task:" + linkMd5
	ok, err := vu.videoRepository.HExists(ctx, key, "task")
	if err != nil {
		return "", err
	}
	if ok {
		return "", errors.New("任务列表存在")
	}
	table := urlutil.GetTable(body.Url)
	if table == "" {
		return "", fetcherr.TableNoLink
	}
	Pan := urlutil.GetPanTable(table, vu.env)
	if Pan == "" {
		return "", errors.New("获取Pan错误")
	}
	exist, err := vu.videoRepository.ExistLinkMd5(ctx, table, linkMd5)
	if err != nil {
		return "", err
	}
	if exist {
		return "", errors.New("数据库存在")
	}
	_, err = vu.videoRepository.HMSet(ctx, key, "task", "1")
	if err != nil {
		return "", err
	}
	lastSlashIndex := strings.LastIndex(body.Surl, "/")
	if lastSlashIndex == -1 {
		return "", errors.New("获取/路径错误")
	}

	baseURL := body.Surl[:lastSlashIndex]
	Data := &domain.RmqPostData{
		Url:       body.Url,
		Surl:      body.Surl,
		TaskId:    body.TaskId,
		BaseUrl:   baseURL,
		LinkMd5:   linkMd5,
		FileName:  key + ".m3u8",
		PathIds:   Pan,
		TableName: table,
	}
	post, err := json.Marshal(Data)
	if err != nil {
		return "", err
	}
	repStr, err := vu.producer.SendMsg(ctx, vu.env.ToPic, "task", post)
	if err != nil {
		return "", err
	}
	return repStr, nil
}

func (vu *videoUsecase) TaskM4(c context.Context, body *domain.TaskM4Body) error {
	ctx, cancel := context.WithTimeout(c, vu.contextTimeout)
	defer cancel()
	redKey := urlutil.EncodeMD5WithSalt(body.Url, []byte(vu.env.CacheKey))
	if body.Time != 0 {
		err := vu.videoRepository.HSetTimeScript(ctx, redKey, body.Surl, "mp4", "1", body.Time)
		if err != nil {
			return err
		}
		return nil
	}
	err := vu.videoRepository.HSetTimeScript(ctx, redKey, body.Surl, "mp4", "0", 7200)
	if err != nil {
		return err
	}
	linkMd5 := urlutil.EncodeMD5(body.Url)
	table := urlutil.GetTable(body.Url)
	if table == "" {
		return fetcherr.TableNoLink
	}
	data := &domain.UPDate{
		Pan:      6,
		PathId:   "mp4",
		FileId:   "mp4",
		LinkMD5:  linkMd5,
		Link:     body.Url,
		FileName: body.Surl,
	}
	err = vu.videoRepository.UPDataByLinkMd5(ctx, table, data)
	if err != nil {
		return err
	}
	return nil
}
func (vu *videoUsecase) DelM3(c context.Context, body *domain.DelM3Body) error {
	ctx, cancel := context.WithTimeout(c, vu.contextTimeout)
	defer cancel()
	var err error
	redKey := urlutil.EncodeMD5WithSalt(body.Url, []byte(vu.env.CacheKey))
	from, err := vu.videoRepository.HGet(ctx, redKey, "from")
	if err != nil {
		return err
	}
	if from == "0" {
		err = vu.videoRepository.Del(ctx, redKey)
		if body.Save {
			linkMd5 := urlutil.EncodeMD5(body.Url)
			table := urlutil.GetTable(body.Url)
			if table == "" {
				return fetcherr.TableNoLink
			}
			data, err := vu.videoRepository.DelDataByLinkMd5(ctx, table, linkMd5)
			if err != nil {
				return err
			}
			switch data.Pan {
			case 3:
				err = vu.yunPan139.Remove(data.FileId)
			case 2:
			}
			if err != nil {
				return err
			}
		}
	}
	if from == "1" {
		_, err := vu.videoRepository.HDel(ctx, redKey, "task")
		rely, err := vu.videoRepository.HIncrBy(ctx, redKey, "count", 1)
		if err != nil {
			return err
		}
		if rely <= 3 {
			_, err = vu.videoRepository.HMSet(ctx, redKey, "url", "")
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func (vu *videoUsecase) DelM4(c context.Context, url string) error {
	ctx, cancel := context.WithTimeout(c, vu.contextTimeout)
	defer cancel()
	redKey := urlutil.EncodeMD5WithSalt(url, []byte(vu.env.CacheKey))
	from, err := vu.videoRepository.HGet(ctx, redKey, "from")
	if err != nil {
		return err
	}
	if from == "0" {
		if err = vu.videoRepository.Del(ctx, redKey); err != nil {
			return err
		}
		linkMd5 := urlutil.EncodeMD5(url)
		table := urlutil.GetTable(url)
		if table == "" {
			return fetcherr.TableNoLink
		}
		err = vu.videoRepository.DelByLinkMd5(ctx, table, linkMd5)
		if err != nil {
			return err
		}
	}
	if from == "1" {
		rely, err := vu.videoRepository.HIncrBy(ctx, redKey, "count", 1)
		if err != nil {
			return err
		}
		if rely <= 3 {
			_, err = vu.videoRepository.HMSet(ctx, redKey, "url", "")
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func (vu *videoUsecase) SaveM3(c context.Context, url string, data *[]byte) error {
	ctx, cancel := context.WithTimeout(c, vu.contextTimeout*3)
	defer cancel()
	linkMd5 := urlutil.EncodeMD5(url)
	table := urlutil.GetTable(url)
	if table == "" {
		return fetcherr.TableNoLink
	}
	redKey := urlutil.EncodeMD5WithSalt(url, []byte(vu.env.CacheKey))
	Pan := urlutil.GetPanTable(table, vu.env)
	if Pan == "" {
		return errors.New("根据链接获取panPath错误")
	}
	fileName := redKey + ".m3u8"
	fileId, err := vu.yunPan139.Put(Pan, fileName, *data)
	if err != nil {
		_ = feishu.SendMsg("上传云盘出错" + err.Error())
		return err
	}
	upData := &domain.UPDate{
		Pan:      3,
		PathId:   Pan,
		FileId:   fileId,
		LinkMD5:  linkMd5,
		Link:     url,
		FileName: redKey + ".m3u8",
	}
	err = vu.videoRepository.UPDataByLinkMd5(ctx, table, upData)
	if err != nil {
		_ = feishu.SendMsg("保存mysql错误" + err.Error())
		err = vu.yunPan139.Remove(fileId)
		if err != nil {
			_ = feishu.SendMsg("删除云盘资源出错" + err.Error())
			return err
		}
		return err
	}
	return nil
}
