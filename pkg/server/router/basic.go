package router

import (
	"ginboot/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

type BasicRouter struct {
	Router
}

func NewBasicRouter() *BasicRouter {
	return &BasicRouter{}
}

func (sr *BasicRouter) RouterGroups() []*RouteGroup {
	return []*RouteGroup{
		{
			Prefix: "/-",
			SubRoutes: []Route{
				{
					Name:    "a health check, just for monitoring",
					Method:  GET,
					Pattern: "/health",
					HandlerFunc: func(ctx *gin.Context) {
						ctx.String(http.StatusOK, "OK")
					},
				},
				{
					Name:    "get the log level",
					Method:  GET,
					Pattern: "/log/level",
					HandlerFunc: func(ctx *gin.Context) {
						logger.SetLevelHTTP(ctx.Writer, ctx.Request)
					},
				},
				{
					Method:  GET,
					Pattern: "/metrics",
					HandlerFunc: func(ctx *gin.Context) {
						gin.WrapH(promhttp.Handler())
					},
				},
				{
					Method:  GET,
					Pattern: "/swagger/*any",
					HandlerFunc: func(ctx *gin.Context) {
						ginSwagger.WrapHandler(swaggerFiles.Handler)
					},
				},
			},
		},
	}
}
