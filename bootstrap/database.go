package bootstrap

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func NewMysqlDatabase(env *Env) (*sqlx.DB, error) {
	dbHost := env.DBHost
	dbPort := env.DBPort
	dbUser := env.DBUser
	dbPass := env.DBPass
	dbData := env.DBDaTa
	var dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbPort, dbData)
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, err
	}
	// 设置连接池的最大闲置连接数和最大打开连接数
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(5 * time.Minute)
	//database.SetConnMaxIdleTime(4 * time.Minute)

	// 封装为 DB 结构体返回
	return db, nil
}
