package route

import (
	"github.com/gofiber/fiber/v2"
	"time"
	"video_cache/api/controller"
	"video_cache/bootstrap"
	"video_cache/internal/rmq"
	"video_cache/repository"
	"video_cache/usecase"
)

func NewVideoRouter(app *bootstrap.Application, timeout time.Duration, group fiber.Router) {

	vr := repository.NewVideoRepository(app.Rdb, app.DB)
	producer := rmq.NewProducer(app.Rdb, rmq.WithMsgQueueLen(50))
	vc := controller.VideoController{
		VideoUsecase: usecase.NewVideoUsecase(vr, app.Env, producer, app.YunPan139, timeout),
	}
	// 获取播放链接
	group.Post("/url", vc.GetUrl)
	// m3u8 提交任务
	group.Post("/task_m3u8", vc.TaskM3U8)
	// mp4 提交任务
	group.Post("/task_mp4", vc.TaskMP4)
	// 删除m3u8
	group.Post("/m3u8_err", vc.DelM3U8)
	// 删除mp4
	group.Post("/mp4_err", vc.DelMP4)
	// 保存m3u8 文件
	group.Post("/save_m3", vc.SaveM3U8)
}
