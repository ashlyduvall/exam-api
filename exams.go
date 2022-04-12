package main

import "github.com/gin-gonic/gin"
import "net/http"

type exam struct {
	ID             int         `json:"id"`
	Syllabus       syllabus    `json:"syllabus"`
	ExamTagset     exam_tagset `json:"exam_tagset"`
	Questions      []question  `json:"questions"`
	CreateDateTime string      `json:"create_date_time"`
	StartDateTime  string      `json:"start_date_time"`
	EndDateTime    string      `json:"end_date_time"`
}

func BuildExamRoutes(router *gin.Engine) {
	router.GET("/exams/all", GetAllExams)
}

func GetAllExams(ret *gin.Context) {
	e := GetExamById(1)
	ret.JSON(http.StatusOK, e)
}
func GetExamById(id int) exam {
	s := GetSyllabusById(1)
	e := exam{
		ID:             1,
		Syllabus:       s,
		CreateDateTime: "2022-04-12T12:38:11Z",
	}
	etc := GetExamTagsetByExam(e)
	e.ExamTagset = etc
	q := GetQuestionsByExam(e)
	e.Questions = q
	return e
}
