package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/Rainmoom/gin-boot/pkg/server/middleware"
	"github.com/Rainmoom/gin-boot/pkg/server/router"
	"github.com/Rainmoom/gin-boot/pkg/storage/cache"
	"github.com/Rainmoom/gin-boot/pkg/util"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Rainmoom/gin-boot/pkg/conf"
	"github.com/Rainmoom/gin-boot/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/oklog/run"
)

type Server struct {
	Routers []router.Router
	Init    func()
	Use     func(r *gin.Engine) error
}

func (s *Server) Run(configFilePath string) {
	//init config
	conf.ParseConfig(conf.Cfg, configFilePath)

	//init logger
	ctx, logCancel := context.WithCancel(context.Background())
	if err := logger.Init(ctx, conf.Cfg.Log); err != nil {
		fmt.Println("err init failed", err)
		os.Exit(1)
	}

	//init local cache
	cache.InitLC()

	//init customize function
	if s.Init != nil {
		s.Init()
	}

	//create goroutine group
	gr := run.Group{}

	{
		//listen the system interrupt call
		term := make(chan os.Signal, 1)
		signal.Notify(term, os.Interrupt, syscall.SIGTERM)
		gr.Add(
			func() error {
				<-term
				logger.Out.Warn("Received SIGTERM, exiting gracefully...")
				return nil
			},
			func(err error) {},
		)
	}
	{
		//add server goroutine to group
		cancel := make(chan struct{})
		gr.Add(func() error {
			//set gin run mode
			gin.SetMode(conf.Cfg.Mode)
			//create gin server
			srv := s.setupServer(conf.Cfg)
			//block goroutine and listen exit call
			s.gracefulExit(srv, cancel)
			return nil
		}, func(err error) {
			close(cancel)
		})
	}

	//run the goroutine group and block
	if err := gr.Run(); err != nil {
		logger.Out.Error(err.Error())
	}

	//exit server and close logger
	logger.Out.Info("exiting")
	logCancel()
}

func (s *Server) setupServer(cfg *conf.ConfigYaml) *http.Server {

	//create http server
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Http.Host, cfg.Http.Port),
		Handler: s.setupRouter(),
	}

	//create server listen goroutine and start
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Out.Error(err.Error())
			os.Exit(1)
		}
	}()

	logger.Out.Info(fmt.Sprintf("start on server:%s", srv.Addr))
	return srv
}

// stop server and gracefulExit
func (s *Server) gracefulExit(srv *http.Server, ch chan struct{}) {
	<-ch

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Out.Error(err.Error())
	}
	logger.Out.Info("Shutdown server ...")
}

func (s *Server) setupRouter() *gin.Engine {
	//create gin handler
	r := gin.New()

	//init customize middleware
	r.Use(
		middleware.GinLogger(),
		middleware.GinCatchError(),
	)

	//init customize middleware
	if s.Use != nil {
		if err := s.Use(r); err != nil {
			logger.Out.Error("Use middleware failed: " + err.Error())
			os.Exit(1)
		}
	}

	{
		//init gin route
		var routeGroup []*router.RouteGroup
		routeGroup = append(routeGroup, router.NewBasicRouter().RouterGroups()...)
		for _, rt := range s.Routers {
			routeGroup = append(routeGroup, rt.RouterGroups()...)
		}

		routeGroupsMap := make(map[string]*gin.RouterGroup)
		for _, gRoute := range routeGroup {
			if _, ok := routeGroupsMap[gRoute.Prefix]; !ok {
				routeGroupsMap[gRoute.Prefix] = r.Group(gRoute.Prefix)
			}
			for _, gMiddleware := range gRoute.GroupMiddleware {
				routeGroupsMap[gRoute.Prefix].Use(gMiddleware)
			}

			for _, subRoute := range gRoute.SubRoutes {
				length := len(subRoute.Middleware) + 2
				routes := make([]any, length)
				routes[0] = subRoute.Pattern
				for i, v := range subRoute.Middleware {
					routes[i+1] = v
				}
				routes[length-1] = subRoute.HandlerFunc

				util.CallReflect(
					routeGroupsMap[gRoute.Prefix],
					subRoute.Method,
					routes...)
			}
		}
	}

	return r
}
