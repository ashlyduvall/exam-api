package main

import "github.com/gin-gonic/gin"
import "net/http"
import "fmt"

func getFoo(ret *gin.Context) {
	ret.JSON(http.StatusOK, gin.H{
		"message": "It works!",
	})
}

func main() {
	listenAddress := "0.0.0.0:8081"

	router := gin.Default()
	router.GET("/", getFoo)

	fmt.Printf("Application listening on %s.\n", listenAddress)
	router.Run(listenAddress)
}
