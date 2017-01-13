package api

import (
	"github.com/gin-gonic/gin"
	"github.com/globocom/grou-loader/loader"
	"os"
)

func NewEngine(port string) *Engine {
	osHostname, _ := os.Hostname()
	return &Engine{engine: gin.Default(), port: port, handler: NewEngineHandler(osHostname, port)}
}

type Engine struct {
	tools   *loader.Tools
	engine  *gin.Engine
	port    string
	handler *EngineHandler
}

func (a *Engine) loadRoutes() {
	v0 := a.engine.Group("/v0")
	v0.GET("/", a.handler.Index)
	v0.GET("/tools", a.handler.Tools)
	v0.GET("/tools/:tool", a.handler.ToolDenominated)
	v0.POST("/tools/:tool", a.handler.Run)
}

func (a *Engine) Start() {
	a.loadRoutes()
	a.engine.Run(":" + a.port)
}
