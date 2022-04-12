package main

import "github.com/gin-gonic/gin"
import "fmt"

func main() {
	listenAddress := "0.0.0.0:8081"

	router := gin.Default()
	BuildTagRoutes(router)

	fmt.Printf("Application listening on %s.\n", listenAddress)
	router.Run(listenAddress)
}
