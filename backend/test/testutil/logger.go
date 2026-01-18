package testutil

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"

	"piemdm/pkg/log"
)

// getLogPath finds the directory of this file, traverses up to find the project root
// (the directory containing go.mod), and returns the absolute path to the 'logs' subdirectory
// within that root. This ensures a consistent log path regardless of where `go test` is run.
func getLogPath() string {
	// Get the path of the current file.
	_, b, _, ok := runtime.Caller(0)
	if !ok {
		// Fallback if caller information is not available.
		return "./logs"
	}

	// Start searching from the directory of the current file.
	dir := filepath.Dir(b)

	// Traverse upwards until we find a directory containing go.mod.
	for {
		// Check if go.mod exists in the current directory.
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			// Found the project root. Return the path to the 'logs' directory inside it.
			return filepath.Join(dir, "logs")
		}

		// Move up to the parent directory.
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached the filesystem root without finding go.mod.
			// This is unexpected; fallback to a relative path.
			return "./logs"
		}
		dir = parent
	}
}

// CreateTestLogger 创建测试用的logger
// 优先使用配置文件，如果配置文件不存在则使用默认配置
func CreateTestLogger() *log.Logger {
	conf := viper.New()

	// 设置默认配置
	conf.SetDefault("log.log_path", getLogPath()) // Use robust, dynamic log path
	conf.SetDefault("log.log_file_name", "test.log")
	conf.SetDefault("log.log_level", "info")
	conf.SetDefault("log.encoding", "console")
	conf.SetDefault("env", "test")

	// 尝试读取配置文件，但不强制要求
	configPath := "../../config/local.yml"
	conf.SetConfigFile(configPath)

	// 如果配置文件存在就读取，不存在就忽略
	if err := conf.ReadInConfig(); err != nil {
		// 配置文件不存在，使用默认配置
		// 不panic，让测试继续运行
	} else {
		// 配置文件存在，使用配置文件中的值
		// 配置文件的值会覆盖默认值
	}

	return log.NewLog(conf)
}
