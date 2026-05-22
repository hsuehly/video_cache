package rmq

//func Setup(app *bootstrap.Application) (func(), error) {
//rdb := app.Rdb
//client := httpclient.NewHttpClient()
//callbackFunc := func(ctx context.Context, msg *MsgEntity) error {
//	//
//	postData := new(domain.RmqPostData)
//	err := json.Unmarshal([]byte(msg.Val), postData)
//	if err != nil {
//		return err
//	}
//	resp, err := client.Do("GET", postData.Surl, nil, true)
//	if err != nil {
//		return err
//	}
//	if resp.StatusCode != 200 {
//		logger.Warn("资源错误")
//		return fmt.Errorf("资源错误：%d", resp.StatusCode)
//	}
//	var buf bytes.Buffer
//	_, err = buf.ReadFrom(resp.Body)
//	resp.Body.Close()
//	var m3u8Data []byte
//	switch postData.TaskId {
//	case 1:
//		m3u8Data, err = m3task.ParseZuoAn(buf, postData.BaseUrl)
//	}
//	if m3u8Data == nil {
//		logger.Warn("m3u8 data nil")
//		return fmt.Errorf("还原出错：%s", err)
//	}
//	fmt.Println(len(m3u8Data), "m3u8Data")
//	//fileId, err := app.YunPan139.Put(postData.PathIds, postData.FileName, m3u8Data)
//	//data := &domain.UPDate{
//	//	Pan:      3,
//	//	PathId:   postData.PathIds,
//	//	FileId:   fileId,
//	//	LinkMD5:  postData.LinkMd5,
//	//	Link:     postData.Url,
//	//	FileName: postData.FileName,
//	//}
//	//_, errs := app.DB.NamedExec("INSERT INTO "+postData.TableName+" (pan, link, link_md5, path_id, file_id ,file_name) VALUES (:pan, :link, :link_md5, :path_id, :file_id , :file_name) ON DUPLICATE KEY UPDATE pan = :pan, path_id = :path_id, file_id = :file_id, file_name = :file_name, status = 0", data)
//	//if errs != nil {
//	//	return errs
//	//}
//
//	fmt.Printf("receive msg, msg id: %s, msg key: %s, msg val: %s \n", msg.MsgID, msg.Key, msg.Val)
//	//fmt.Printf("url %s, link %s, Surl %s \n", postData.Url, postData.LinkMd5, postData.Surl)
//	return nil
//}
//consumer, err := NewConsumer(rdb, app.Env.ToPic, app.Env.ConsumerGroup, app.Env.ConsumerID, callbackFunc,
//	// 每条消息最多重试 2 次
//	WithMaxRetryLimit(2),
//	// 每轮接收消息的超时时间为 2 s
//	WithReceiveTimeout(2*time.Second))
//if err != nil {
//	return nil, fmt.Errorf("new consumer err %s", err)
//}
//
//return consumer.Stop, nil
//}
