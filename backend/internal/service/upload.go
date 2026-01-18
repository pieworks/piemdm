package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type UploadService struct {
	uploadDir string
	maxSize   int64
}

func NewUploadService() *UploadService {
	uploadDir := "./uploads"

	// 确保上传目录存在
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		panic(fmt.Sprintf("创建上传目录失败: %v", err))
	}

	return &UploadService{
		uploadDir: uploadDir,
		maxSize:   10 << 20, // 10 MB 默认限制
	}
}

// ValidateFile 验证文件
func (s *UploadService) ValidateFile(file *multipart.FileHeader) error {
	// 验证文件大小
	if file.Size > s.maxSize {
		return fmt.Errorf("文件大小超过限制: 最大 %d MB", s.maxSize/(1<<20))
	}

	// 验证文件类型
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := []string{
		".jpg", ".jpeg", ".png", ".gif", ".webp", // 图片
		".pdf", ".doc", ".docx", ".xls", ".xlsx", // 文档
		".txt", ".csv", ".zip", ".rar", // 其他
	}

	allowed := false
	for _, allowedExt := range allowedExts {
		if ext == allowedExt {
			allowed = true
			break
		}
	}

	if !allowed {
		return fmt.Errorf("不支持的文件类型: %s", ext)
	}

	return nil
}

// GenerateFilename 生成唯一文件名
func (s *UploadService) GenerateFilename(ext string) string {
	// 使用时间戳 + 随机字符串
	timestamp := time.Now().Format("20060102_150405")

	// 生成随机字符串
	randomBytes := make([]byte, 4)
	rand.Read(randomBytes)
	randomStr := hex.EncodeToString(randomBytes)

	return fmt.Sprintf("%s_%s%s", timestamp, randomStr, ext)
}

// SaveFile 保存文件到本地
func (s *UploadService) SaveFile(c *gin.Context, file *multipart.FileHeader, filename string) (string, error) {
	// 构建完整路径
	fullPath := filepath.Join(s.uploadDir, filename)

	// 保存文件
	if err := c.SaveUploadedFile(file, fullPath); err != nil {
		return "", err
	}

	// 返回相对路径 (用于存储到数据库)
	return "/uploads/" + filename, nil
}
