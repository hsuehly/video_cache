package controller

import (
	"database/sql"
	"errors"
	"github.com/gofiber/fiber/v2"
	"io"
	"video_cache/domain"
	"video_cache/internal/fetcherr"
	"video_cache/pkg/logger"
)

type VideoController struct {
	VideoUsecase domain.VideoUsecase
}

func (vc *VideoController) GetUrl(c *fiber.Ctx) error {

	body := new(domain.UrlBody)
	if err := c.BodyParser(body); err != nil {
		return c.JSON(fiber.Map{
			"code": 500,
			"msg":  "参数错误",
		})
	}
	data, err := vc.VideoUsecase.GetLink(c.Context(), body.Url)
	if err != nil {
		// 排除查询为空的打印
		if !errors.Is(err, sql.ErrNoRows) && !errors.Is(err, fetcherr.TableNoLink) {
			logger.Warn("获取地址错误", logger.String("url", body.Url), logger.Err(err))
		}
		return c.JSON(fiber.Map{
			"code": 500,
			"msg":  err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"code": 200,
		"data": data,
	})
}
func (vc *VideoController) TaskM3U8(c *fiber.Ctx) error {
	body := new(domain.TaskM3Body)
	if err := c.BodyParser(body); err != nil {
		return c.JSON(fiber.Map{
			"code": 500,
			"msg":  err.Error(),
		})
	}
	if body.Url == "" || body.Surl == "" {
		return c.JSON(fiber.Map{
			"code": 500,
			"msg":  "参数错误",
		})

	}
	data, err := vc.VideoUsecase.TaskM3(c.Context(), body)
	if err != nil {
		logger.Warn("提交m3任务失败", logger.String("url", body.Url), logger.String("surl", body.Surl), logger.Err(err))
		return c.JSON(fiber.Map{
			"code": 500,
			"msg":  err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"code": 200,
		"data": data,
	})
}
func (vc *VideoController) TaskMP4(c *fiber.Ctx) error {
	body := new(domain.TaskM4Body)
	if err := c.BodyParser(body); err != nil {
		return c.JSON(fiber.Map{
			"code": 500,
			"msg":  err.Error(),
		})
	}
	if body.Url == "" || body.Surl == "" {
		return c.JSON(fiber.Map{
			"code": 500,
			"msg":  "参数错误",
		})
	}
	err := vc.VideoUsecase.TaskM4(c.Context(), body)
	if err != nil {
		logger.Warn("提交m4任务失败", logger.String("url", body.Url), logger.String("surl", body.Surl), logger.Err(err))
		return c.JSON(fiber.Map{
			"code": 500,
			"msg":  err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"code": 200,
		"data": "ok",
	})
}
func (vc *VideoController) DelM3U8(c *fiber.Ctx) error {
	body := new(domain.DelM3Body)
	if err := c.BodyParser(body); err != nil {
		return c.JSON(fiber.Map{
			"code": 500,
			"msg":  err.Error(),
		})
	}
	err := vc.VideoUsecase.DelM3(c.Context(), body)
	if err != nil {
		logger.Warn("删除m3任务失败", logger.String("url", body.Url), logger.Err(err))
		return c.JSON(fiber.Map{
			"code": 500,
			"msg":  err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"code": 200,
		"data": "ok",
	})
}
func (vc *VideoController) DelMP4(c *fiber.Ctx) error {
	body := new(domain.DelM4Body)
	if err := c.BodyParser(body); err != nil {
		return c.JSON(fiber.Map{
			"code": 500,
			"msg":  err.Error(),
		})
	}
	if body.Url == "" {
		return c.JSON(fiber.Map{
			"code": 500,
			"msg":  "参数错误",
		})
	}
	err := vc.VideoUsecase.DelM4(c.Context(), body.Url)
	if err != nil {
		logger.Warn("删除m4任务失败", logger.String("url", body.Url), logger.Err(err))
		return c.JSON(fiber.Map{
			"code": 500,
			"msg":  err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"code": 200,
		"data": "ok",
	})
}
func (vc *VideoController) SaveM3U8(c *fiber.Ctx) error {
	url := c.FormValue("url")
	if url == "" {
		return c.JSON(fiber.Map{
			"code": 500,
			"msg":  "无链接",
		})
	}
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(fiber.Map{
			"code": 500,
			"msg":  "无文件",
		})
	}

	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		return c.JSON(fiber.Map{
			"code": 500,
			"msg":  "文件读取错误",
		})
	}
	defer src.Close()

	// 将上传的文件内容复制到目标文件
	// 读取上传的文件内容
	data, err := io.ReadAll(src)
	if err != nil {
		return c.JSON(fiber.Map{
			"code": 500,
			"msg":  "文件读取错误2",
		})
	}
	if len(data) == 0 {
		return c.JSON(fiber.Map{
			"code": 500,
			"msg":  "空文件",
		})

	}
	err = vc.VideoUsecase.SaveM3(c.Context(), url, &data)
	if err != nil {
		logger.Warn("保存m3任务失败", logger.String("url", url), logger.Err(err))
		return c.JSON(fiber.Map{
			"code": 500,
			"msg":  err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"code": 200,
		"data": "ok",
	})
}
