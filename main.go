package main

import (
	"flag"
	"github.com/globocom/grou-loader/api"
)

var (
	port = flag.String("PORT", "4082", "Server default port")
)

func main() {
	flag.Parse()
	engine := api.NewEngine(*port)
	engine.Start()
}
