package main

import (
	"fmt"
	"time"
	"video_cache/internal/rmq"
	"video_cache/pkg/logger"

	"video_cache/api/route"
	"video_cache/bootstrap"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app, err := bootstrap.App()
	if err != nil {
		panic(fmt.Sprintf("bootstrap app init err %s", err))
	}
	lg, err := logger.Init(logger.WithEnv(app.Env.AppEnv))
	if err != nil {
		panic(fmt.Sprintf("logger init err %s", err))
	}
	f := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})
	route.Setup(app, f)
	cons, err := rmq.NewConsumer(app, rmq.WithReceiveTimeout(3*time.Second))
	if err != nil {
		panic(fmt.Sprintf("rmq Setup init err %s", err))
	}
	defer func() {
		app.CloseDBConnection()
		app.CloseRdbConnection()
		app.CloseYun139Connection()
		cons.Stop()
		lg.Sync()
	}()
	if err = f.Listen(app.Env.PORT); err != nil {
		panic(fmt.Sprintf("fiber listen err %s", err))

	}
}
