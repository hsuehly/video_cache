package repository

import (
	"context"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"time"
	"video_cache/domain"
	"video_cache/internal/redisscript"
)

type videoRepository struct {
	rdb      *redis.Client
	database *sqlx.DB
}

func NewVideoRepository(rdb *redis.Client, db *sqlx.DB) domain.VideoRepository {
	return &videoRepository{
		rdb:      rdb,
		database: db,
	}
}

// Redis

func (vr *videoRepository) Exists(c context.Context, keys ...string) (int64, error) {
	return vr.rdb.Exists(c, keys...).Result()
}
func (vr *videoRepository) SAdd(c context.Context, key string, members ...interface{}) (int64, error) {
	return vr.rdb.SAdd(c, key, members).Result()
}
func (vr *videoRepository) SRem(c context.Context, key string, members ...interface{}) (int64, error) {
	return vr.rdb.SRem(c, key, members).Result()
}
func (vr *videoRepository) HGetAll(c context.Context, key string) (map[string]string, error) {
	return vr.rdb.HGetAll(c, key).Result()
}
func (vr *videoRepository) HGetAllTime(c context.Context, key string) (map[string]string, int, error) {
	result, err := redisscript.HGetAllTimeScript.Run(c, vr.rdb, []string{key}).Result()
	if err != nil {
		return nil, 0, err
	}
	data, ok := result.([]any)
	if !ok {
		return nil, 0, errors.New("result 断言错误")
	}
	if len(data) != 2 {
		return nil, 0, errors.New("result data 错误")
	}
	ttl, ok := data[1].(int64)
	if !ok {
		return nil, 0, errors.New("data[1] 断言错误")
	}
	info, ok := data[0].([]any)
	if !ok {
		return nil, 0, errors.New("data[0] 断言错误")
	}
	hashMap := make(map[string]string)
	for i := 0; i < len(info); i += 2 {
		field := info[i].(string)
		value := info[i+1].(string)
		hashMap[field] = value
	}
	return hashMap, int(ttl), nil
}
func (vr *videoRepository) HGetUrlTypeTime(c context.Context, key string) (Url, Type string, TTl int) {
	result := redisscript.HGetUrlTypeTimeScript.Run(c, vr.rdb, []string{key}, "url", "type").Val()
	data, ok := result.([]any)
	if !ok || len(data) != 2 {
		return
	}
	ttl, ok := data[1].(int64)
	if !ok || ttl <= 0 {
		return
	}
	info, ok := data[0].([]any)
	if !ok {
		return
	}
	Url, ok = info[0].(string)
	if !ok {
		return
	}
	Type, ok = info[1].(string)
	if !ok {
		return
	}
	TTl = int(ttl)
	return

}
func (vr *videoRepository) HSetTimeScript(c context.Context, key, url, urlType, from string, expTime int) error {
	return redisscript.HSetTimeScript.Run(c, vr.rdb, []string{key}, url, urlType, from, expTime).Err()
}
func (vr *videoRepository) SetNX(c context.Context, key string, value any, seconds time.Duration) (bool, error) {
	return vr.rdb.SetNX(c, key, value, seconds).Result()
}
func (vr *videoRepository) Exist(c context.Context, key string) bool {
	repl := vr.rdb.Exists(c, key).Val()
	if repl == 1 {
		return true
	}
	return false
}
func (vr *videoRepository) Del(c context.Context, key string) error {
	return vr.rdb.Del(c, key).Err()
}
func (vr *videoRepository) HGet(c context.Context, key, field string) (string, error) {
	return vr.rdb.HGet(c, key, field).Result()
}
func (vr *videoRepository) HSetNX(c context.Context, key, field, value string) (bool, error) {
	return vr.rdb.HSetNX(c, key, field, value).Result()
}
func (vr *videoRepository) HDel(c context.Context, key string, field ...string) (int64, error) {
	return vr.rdb.HDel(c, key, field...).Result()
}
func (vr *videoRepository) HIncrBy(c context.Context, key string, field string, incr int64) (int64, error) {
	return vr.rdb.HIncrBy(c, key, field, incr).Result()
}
func (vr *videoRepository) HMSet(c context.Context, key string, values ...any) (bool, error) {
	return vr.rdb.HMSet(c, key, values).Result()
}
func (vr *videoRepository) HExists(c context.Context, key string, field string) (bool, error) {
	return vr.rdb.HExists(c, key, field).Result()
}

// Mysql

func (vr *videoRepository) GetLinkByMd5(c context.Context, tableName, linkMd5 string) (*domain.HlsLink, error) {
	links := new(domain.HlsLink)
	selectLink := "SELECT pan,file_id,file_name FROM " + tableName + "  WHERE link_md5 = ? AND status = 0"
	err := vr.database.GetContext(c, links, selectLink, linkMd5)
	if err != nil {
		return links, err
	}
	return links, nil
}

// mp4Err

func (vr *videoRepository) DelByLinkMd5(c context.Context, tableName, linkMd5 string) error {
	upDataSql := "UPDATE " + tableName + " SET status = 1 WHERE link_md5 = ?"
	_, err := vr.database.ExecContext(c, upDataSql, linkMd5)
	if err != nil {
		return err
	}
	return nil
}

func (vr *videoRepository) ExistLinkMd5(c context.Context, tableName, linkMd5 string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM " + tableName + " WHERE link_md5 = ? AND status = 0)"
	err := vr.database.GetContext(c, &exists, query, linkMd5)
	if err != nil {
		return false, err
	}
	return exists, nil

}

func (vr *videoRepository) UPDataByLinkMd5(c context.Context, table string, data *domain.UPDate) error {
	upDataSql := "INSERT INTO " + table + " (pan, link, link_md5, path_id, file_id ,file_name) VALUES (:pan, :link, :link_md5, :path_id, :file_id , :file_name) ON DUPLICATE KEY UPDATE pan = :pan, path_id = :path_id, file_id = :file_id, file_name = :file_name, status = 0"
	_, err := vr.database.NamedExecContext(c, upDataSql, data)
	if err != nil {
		return err
	}
	return nil
}
func (vr *videoRepository) DelDataByLinkMd5(c context.Context, table, linkMd5 string) (*domain.PanRecord, error) {
	upDataSql := "UPDATE " + table + " SET status = 1 WHERE link_md5 = ?"
	selectQuery := "SELECT pan, file_id, file_name FROM  " + table + "  WHERE link_md5 = ?"
	_, err := vr.database.ExecContext(c, upDataSql, linkMd5)
	if err != nil {
		return nil, err
	}
	var data = new(domain.PanRecord)
	err = vr.database.GetContext(c, data, selectQuery, linkMd5)
	if err != nil {
		return nil, err
	}
	return data, nil
}
