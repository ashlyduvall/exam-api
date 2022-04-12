package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type exam_tagset struct {
	ID          int       `json:"id"`
	Syllabus    *syllabus `json:"syllabus"`
	DisplayName string    `json:"display_name"`
	Tags        []tag     `json:"tags"`
}

func BuildExamTagsetRoutes(router *gin.Engine) {
	router.GET("/exam_tagsets/all", GetAllExamTagsets)
}

func GetAllExamTagsets(ret *gin.Context) {
	e, err := GetExamById(1)

	if err != nil {
		fmt.Println("Error getting exam tagsets!")
		fmt.Println(err)
		ret.JSON(http.StatusInternalServerError, gin.H{
			"Message": "Error getting exam tagsets!",
		})
		return
	}

	ets, err := GetExamTagsetByExam(*e)

	if err != nil {
		fmt.Println("Error getting exam tagsets for this exam!")
		fmt.Println(e)
		fmt.Println(err)
		ret.JSON(http.StatusInternalServerError, gin.H{
			"Message": "Error getting exam tagsets!",
		})
		return
	}

	ret.JSON(http.StatusOK, ets)
}

func GetExamTagsetByExam(e exam) (*exam_tagset, error) {
	s, err := GetSyllabusById(1)

	if err != nil {
		return nil, err
	}

	t := []tag{
		{
			ID:          1,
			Syllabus:    s,
			DisplayName: "Some Tag",
		},
	}
	return &exam_tagset{
		ID:          1,
		Syllabus:    s,
		DisplayName: "Some exam tagset",
		Tags:        t,
	}, nil
}
