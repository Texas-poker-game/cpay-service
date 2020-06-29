package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"queding.com/go/common/config"
	"queding.com/go/common/log"

	"cpay/eos/api"
)

func setupRouter() *gin.Engine {
	r := gin.New()
	p := ginprometheus.NewPrometheus("cpay")
	p.Use(r)

	accessLogFile := log.AccessLogFile()
	logger := log.AppLogger()

	r.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: func(param gin.LogFormatterParams) string {
			// your custom format
			return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s\" \"%s\" %s\n",
				param.ClientIP,
				param.TimeStamp.Format(time.RFC1123),
				param.Method,
				param.Path,
				param.Request.Proto,
				param.StatusCode,
				param.Latency,
				param.Request.UserAgent(),
				param.ErrorMessage,
			)
		},
		Output: accessLogFile,
	}))

	// recover
	r.Use(func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logger.Error(r)
				result := gin.H{"success": false, "code": 500, "message": "未知错误"}
				c.JSON(http.StatusOK, result)
			}
		}()
		c.Next()
		if err := c.Errors.Last(); err != nil {
			result := gin.H{"success": false, "code": 400, "message": err.Error()}
			c.JSON(http.StatusOK, result)
		}
	})

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// 提现
	r.POST("/eos/withdraw", api.WithdrawHandler)

	// auth token
	r.POST("/eos/auth/token", api.NewAuthToken)

	return r
}

func Start() {
	logger := log.AppLogger()
	address := config.GetString("server.api.address")
	logger.Infof("server will listen on %s", address)

	if err := setupRouter().Run(address); err != nil {
		panic(err)
	}
}
