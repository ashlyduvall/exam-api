package main

import "github.com/gin-gonic/gin"
import "net/http"

type tag struct {
	ID          int `json:"id"`
	Syllabus    syllabus
	DisplayName string `json:"display_name"`
}

func BuildTagRoutes(router *gin.Engine) {
	router.GET("/tags/get/all", GetAllTags)
}

func GetAllTags(ret *gin.Context) {
	s := GetSyllabusById(1)
	t := tag{
		ID:          1,
		Syllabus:    s,
		DisplayName: "Some Tag",
	}

	ret.JSON(http.StatusOK, t)
}
func GetTagsByQuestion(q question) []tag {
	s := GetSyllabusById(1)
	return []tag{
		{ID: 1, Syllabus: s, DisplayName: "Some Tag"},
	}
}
