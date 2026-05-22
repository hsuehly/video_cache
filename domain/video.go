package domain

import (
	"context"
	"time"
)

type Video struct {
	Link    string `db:"link"`
	LinkMd5 string `db:"link_md5"`
	FileID  string `db:"file_id"`
}
type HlsLink struct {
	Pan int `db:"pan"`
	//PathId   string `db:"path_id"`
	FileId   string `db:"file_id"`
	FileName string `db:"file_name"`
}
type UPDate struct {
	Pan      int    `db:"pan"`
	PathId   string `db:"path_id"`
	FileId   string `db:"file_id"`
	LinkMD5  string `db:"link_md5"`
	Link     string `db:"link"`
	FileName string `db:"file_name"`
}
type PanRecord struct {
	Pan      int    `db:"pan"`
	FileId   string `db:"file_id"`
	FileName string `db:"file_name"`
}
type VideoRepository interface {
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	HGetAllTime(ctx context.Context, key string) (map[string]string, int, error)
	HSetTimeScript(ctx context.Context, key, url, urlType, from string, expTime int) error
	HGetUrlTypeTime(c context.Context, key string) (Url, Type string, TTl int)
	Del(ctx context.Context, key string) error
	Exists(ctx context.Context, keys ...string) (int64, error)
	SAdd(ctx context.Context, key string, members ...interface{}) (int64, error)
	SRem(ctx context.Context, key string, members ...interface{}) (int64, error)
	SetNX(c context.Context, key string, value any, seconds time.Duration) (bool, error)
	Exist(c context.Context, key string) bool
	HGet(c context.Context, key, field string) (string, error)
	HSetNX(c context.Context, key, field, value string) (bool, error)
	HDel(c context.Context, key string, field ...string) (int64, error)
	HIncrBy(c context.Context, key string, field string, incr int64) (int64, error)
	HMSet(c context.Context, key string, values ...any) (bool, error)
	HExists(c context.Context, key string, field string) (bool, error)

	// mysql

	GetLinkByMd5(ctx context.Context, tableName, linkMd5 string) (*HlsLink, error)
	DelByLinkMd5(ctx context.Context, tableName, linkMd5 string) error
	ExistLinkMd5(ctx context.Context, tableName, linkMd5 string) (bool, error)
	UPDataByLinkMd5(ctx context.Context, table string, data *UPDate) error
	DelDataByLinkMd5(ctx context.Context, table string, linkMd5 string) (*PanRecord, error)
}
