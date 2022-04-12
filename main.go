package main

import "github.com/gin-gonic/gin"
import "os"

func getConfig(key string, default_value string) string {
	value, found := os.LookupEnv(key)
	if found {
		return value
	} else {
		return default_value
	}
}

func main() {
	listenAddress := getConfig("LISTEN_ADDRESS", "0.0.0.0:8081")

	router := gin.Default()
	BuildTagRoutes(router)
	BuildQuestionRoutes(router)
	BuildQuestionAnswerRoutes(router)
	BuildExamRoutes(router)
	BuildExamTagsetRoutes(router)

	router.Run(listenAddress)
}
