package router

import (
	"github.com/gin-gonic/gin"
)

var (
	GET     = "GET"
	POST    = "POST"
	PUT     = "PUT"
	DELETE  = "DELETE"
	PATCH   = "PATCH"
	UPDATE  = "UPDATE"
	OPTIONS = "OPTIONS"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc func(ctx *gin.Context)
	Middleware  []gin.HandlerFunc
}

type RouteGroup struct {
	Prefix          string
	GroupMiddleware gin.HandlersChain
	SubRoutes       []Route
}

type Router interface {
	RouterGroups() []*RouteGroup
}
