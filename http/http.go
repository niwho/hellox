package http

import (
	"context"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/niwho/hellox/config"
	"github.com/niwho/hellox/logs"
)

// 注意在main函数的末尾调用，会阻塞
func InitHttp() {
	gin.SetMode(config.Conf.Core.Mode)
	router := gin.Default()
	pprof.Register(router, nil)

	router.Static("/static", "./static/static")
	router.StaticFile("/index", "./static/index.html")
	router.StaticFile("/", "./static/index.html")
	router.StaticFile("/favicon.ico", "./static/favicon.ico")
	router.StaticFile("/manifest.json", "./static/manifest.json")

	router.POST("/stat", stat)
	router.GET("/ws", serveWS)
	router.GET("/broadcast", broadcast)
	// router .Run(config.Conf.API.BindAddr)
	srv := &http.Server{
		Addr:    config.Conf.Core.BindAddr,
		Handler: router,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil {
			logs.Log(nil).Info("listen: %s\n", err)
		}
	}()
	logs.Log(logs.F{"listen addr": config.Conf.Core.BindAddr}).Info("service started")
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logs.Log(nil).Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logs.Log(nil).Fatal("Server Shutdown:", err)
	}
	logs.Log(nil).Info("Server exiting")
}
