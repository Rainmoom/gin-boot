package main

import (
	"gin-boot/pkg/conf"
	"gin-boot/pkg/logger"
	"gin-boot/pkg/server"
	"gin-boot/pkg/server/router"
	"gin-boot/pkg/storage/cache"
	"gin-boot/pkg/storage/db/postgres"
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
