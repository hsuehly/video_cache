package route

import (
	"runtime/debug"
	"time"
	"video_cache/pkg/logger"

	"video_cache/bootstrap"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func Setup(app *bootstrap.Application, f *fiber.App) {
	timeout := time.Duration(app.Env.ContextTimeout) * time.Second
	f.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			logger.Error("Recovered from panic",
				logger.String("method", c.Method()),
				logger.String("path", c.Path()),
				logger.Any("panic", e),
				logger.ByteString("stack", debug.Stack()),
			)
		},
	}))
	apiV1 := f.Group("/api/v1")
	NewVideoRouter(app, timeout, apiV1)

}
