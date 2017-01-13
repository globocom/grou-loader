package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/globocom/grou-loader/loader"
)

func NewEngineHandler(hostname string, port string) *EngineHandler {
	tools := loader.NewTools()
	tools.Add(loader.NewVegeta())
	return &EngineHandler{tools: tools, hostname: hostname, port: port}
}

type EngineHandler struct {
	hostname string
	port     string
	tools    *loader.Tools
}

// Index
func (e *EngineHandler) Index(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"tools": href(e.hostname, e.port, "")})
}

// Tools ...
func (e *EngineHandler) Tools(ctx *gin.Context) {
	tools := make([]interface{}, 0)
	for _, tool := range e.tools.All() {
		if tool != nil {
			toolName := tool.GetName()
			href := href(e.hostname, e.port, toolName)
			status := status(tool.Check())
			tools = append(tools, map[string]interface{}{"name": toolName, "status": status, "conf": tool.Params(), "href": href})
		}

	}
	ctx.JSON(200, gin.H{"tools": tools})
}

// Tool Denominated ...
func (e *EngineHandler) ToolDenominated(ctx *gin.Context) {
	toolName := ctx.Param("tool")

	tool := e.tools.FindByName(toolName)
	if tool != nil {
		href := href(e.hostname, e.port, toolName)
		status := status(tool.Check())
		toolsJson := map[string]interface{}{"name": toolName, "status": status, "conf": tool.Params(), "href": href}
		ctx.JSON(200, gin.H{"tool": toolsJson})
	} else {
		ctx.JSON(404, gin.H{"error": "Not Found"})
	}
}

// Run ...
func (e *EngineHandler) Run(ctx *gin.Context) {
	toolName := ctx.Param("tool")
	params := processBody(ctx)

	tool := e.tools.FindByName(toolName)
	href := href(e.hostname, e.port, toolName)
	status := "waiting"
	toolsJson := map[string]interface{}{"name": toolName, "status": status, "conf": tool.Params(), "href": href}
	tool.Start(params)
	ctx.JSON(201, gin.H{"tool": toolsJson})
}

func href(hostname, port, name string) string {
	return fmt.Sprintf("http://%s:%s/v0/tools/%s", hostname, port, name)
}

func status(running bool) string {
	return map[bool]string{true: "running", false: "stopped"}[running]
}

func processBody(ctx *gin.Context) map[string]interface{} {
	body := new(bytes.Buffer)
	body.ReadFrom(ctx.Request.Body)
	var params map[string]interface{}
	json.Unmarshal(body.Bytes(), &params)
	return params
}
