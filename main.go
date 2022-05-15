package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

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
	router.Use(CORSMiddleware())
	BuildTagRoutes(router)
	BuildQuestionRoutes(router)
	BuildExamRoutes(router)

	router.Run(listenAddress)
}
