package main

import "github.com/gin-gonic/gin"
import "net/http"

type exam_tagset struct {
	ID          int      `json:"id"`
	Syllabus    syllabus `json:"syllabus"`
	DisplayName string   `json:"display_name"`
	Tags        []tag    `json:"tags"`
}

func BuildExamTagsetRoutes(router *gin.Engine) {
	router.GET("/exam_tagsets/all", GetAllExamTagsets)
}

func GetAllExamTagsets(ret *gin.Context) {
	e := GetExamById(1)
	ets := GetExamTagsetByExam(e)
	ret.JSON(http.StatusOK, ets)
}

func GetExamTagsetByExam(e exam) exam_tagset {
	s := GetSyllabusById(1)
	t := []tag{
		{
			ID:          1,
			Syllabus:    s,
			DisplayName: "Some Tag",
		},
	}
	return exam_tagset{
		ID:          1,
		Syllabus:    s,
		DisplayName: "Some exam tagset",
		Tags:        t,
	}
}
