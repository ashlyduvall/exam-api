package main

import "github.com/gin-gonic/gin"

func main() {
	listenAddress := "0.0.0.0:8081"

	router := gin.Default()
	BuildTagRoutes(router)
	BuildQuestionRoutes(router)
	BuildQuestionAnswerRoutes(router)

	router.Run(listenAddress)
}
