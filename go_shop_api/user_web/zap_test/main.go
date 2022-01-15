package main

import (
	"time"

	"go.uber.org/zap"
)

func main() {

	logger, _ := zap.NewProduction()

	// 调试级别
	// logger, _ := zap.NewDevelopment()

	defer logger.Sync() // flushes buffer, if any

	url := "http://github.com"
	sugar := logger.Sugar()
	sugar.Infow("failed to fetch URL",
		// Structured context as loosely typed key-value pairs.
		"url", url,
		"attempt", 3,
		"backoff", time.Second,
	)
	sugar.Infof("Failed to fetch URL: %s", url)

	logger.Info("failed to fetch URL", zap.String("url", url),
		zap.Int("num", 222))

}
