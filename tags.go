package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type tag struct {
	ID          int `json:"id"`
	Syllabus    *syllabus
	DisplayName string `json:"display_name"`
}

func BuildTagRoutes(router *gin.Engine) {
	router.GET("/tags/all", GetAllTags)
}

func GetAllTags(ret *gin.Context) {
	s, err := GetSyllabusById(1)

	if err != nil {
		fmt.Printf("Error getting syllabus for tag: %v\n", err)
		ret.JSON(http.StatusInternalServerError, err)
		return
	}

	t := tag{
		ID:          1,
		Syllabus:    s,
		DisplayName: "Some Tag",
	}

	ret.JSON(http.StatusOK, t)
}
func GetTagsByQuestion(q question) (*[]*tag, error) {
	s, err := GetSyllabusById(1)

	if err != nil {
		return nil, err
	}

	t := tag{ID: 1, Syllabus: s, DisplayName: "Some Tag"}
	tt := []*tag{&t}

	return &tt, nil
}
