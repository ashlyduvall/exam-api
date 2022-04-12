package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	listenAddress := GetConfig("LISTEN_ADDRESS", "0.0.0.0:8081")

	// Check DB connectivity
	err := ConnectAndTestDB()

	if err != nil {
		fmt.Println("Unable to connect to database!")
		fmt.Println(err)
		return
	}

	defer DB.Close()

	router := gin.Default()
	BuildTagRoutes(router)
	BuildQuestionRoutes(router)
	BuildExamRoutes(router)
	BuildExamTagsetRoutes(router)

	router.Run(listenAddress)
}
