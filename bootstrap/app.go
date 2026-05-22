package bootstrap

import (
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"video_cache/pkg/yun139"
)

type Application struct {
	Env       *Env
	DB        *sqlx.DB
	Rdb       *redis.Client
	YunPan139 *yun139.Yun139
}

func App() (*Application, error) {
	app := &Application{}
	app.Env = NewEnv()
	app.Rdb = NewRedisClient(app.Env)
	db, err := NewMysqlDatabase(app.Env)
	if err != nil {
		return nil, err
	}
	app.DB = db
	yun, err := yun139.New139Yun(app.Env.Yun139)
	if err != nil {
		return nil, err
	}
	app.YunPan139 = yun
	return app, nil
}
func (app *Application) CloseDBConnection() {
	app.DB.Close()
}
func (app *Application) CloseRdbConnection() {
	app.Rdb.Close()
}
func (app *Application) CloseYun139Connection() {
	app.YunPan139.Destroy()
}
