package main

import (
	"flag"

	"github.com/gin-gonic/gin"
	"github.com/theojulienne/go-wireless/api"
)

func main() {
	var bind string
	flag.StringVar(&bind, "b", ":8081", "the address/port to bind to")
	flag.Parse()

	r := gin.Default()
	api.SetupRoutes(r)
	r.Run(bind)
}
