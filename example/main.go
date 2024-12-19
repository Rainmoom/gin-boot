package main

import (
	"github.com/Rainmoom/gin-boot/pkg/conf"
	"github.com/Rainmoom/gin-boot/pkg/logger"
	"github.com/Rainmoom/gin-boot/pkg/server"
	"github.com/Rainmoom/gin-boot/pkg/server/router"
	"github.com/Rainmoom/gin-boot/pkg/storage/cache"
	"github.com/Rainmoom/gin-boot/pkg/storage/db/postgres"
	"os"
)

func main() {
	svr := &server.Server{
		Init: func() {
			//初始化
			if err := postgres.Init(conf.Cfg.Postgres); err != nil {
				logger.Out.Error("postgres init failed: " + err.Error())
				os.Exit(1)
			}

			if err := cache.InitRC(conf.Cfg.Redis); err != nil {
				logger.Out.Error("redis init failed: " + err.Error())
				os.Exit(1)
			}
		},
		Routers: []router.Router{},
	}
	svr.Run("./example/config.yaml")
}
