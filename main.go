package main

import (
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"log/slog"
	"os"
	"time"
)

func main() {
	// ファイルを開く
	file, err := os.OpenFile("server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error")
		return
	}
	defer file.Close()
	fileWriter := io.Writer(file)

	// loggerの生成
	logger := slog.New(slog.NewJSONHandler(fileWriter, nil))
	// 生成したloggerをデフォルトにする
	slog.SetDefault(logger)

	r := gin.Default()
	// Ginのミドルウェアを設定
	r.Use(ginLogFormat(logger))

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	r.Run()
}

// Ginのミドルウェア
func ginLogFormat(logger *slog.Logger) gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		logger.Info("gin-request",
			slog.String("time", param.TimeStamp.Format(time.RFC3339)),
			slog.Int("status", param.StatusCode),
			slog.String("latency", param.Latency.String()),
			slog.String("client_ip", param.ClientIP),
			slog.String("method", param.Method),
			slog.String("path", param.Path),
			slog.String("error", param.ErrorMessage),
		)
		return ""
	})
}
