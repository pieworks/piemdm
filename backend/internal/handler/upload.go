package handler

import (
	"path/filepath"

	"piemdm/internal/service"

	"github.com/gin-gonic/gin"
)

type UploadHandler interface {
	// 业务方法
	Upload(c *gin.Context)
}

type uploadHandler struct {
	*Handler
	uploadService *service.UploadService
}

func NewUploadHandler(handler *Handler, uploadService *service.UploadService) UploadHandler {
	return &uploadHandler{
		Handler:       handler,
		uploadService: uploadService,
	}
}

// Upload 处理文件上传
func (h *uploadHandler) Upload(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "文件上传失败: " + err.Error()})
		return
	}

	// 验证文件
	if err := h.uploadService.ValidateFile(file); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 生成唯一文件名
	ext := filepath.Ext(file.Filename)
	filename := h.uploadService.GenerateFilename(ext)

	// 保存文件
	filepath, err := h.uploadService.SaveFile(c, file, filename)
	if err != nil {
		c.JSON(500, gin.H{"error": "文件保存失败: " + err.Error()})
		return
	}

	// 返回文件信息
	c.JSON(200, gin.H{
		"url":      filepath,
		"filename": file.Filename,
		"size":     file.Size,
	})
}
