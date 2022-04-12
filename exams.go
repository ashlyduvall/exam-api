package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type exam struct {
	ID             int          `json:"id"`
	Syllabus       *syllabus    `json:"syllabus"`
	ExamTagset     *exam_tagset `json:"exam_tagset"`
	Questions      *[]*question `json:"questions"`
	CreateDateTime string       `json:"create_date_time"`
	StartDateTime  string       `json:"start_date_time"`
	EndDateTime    string       `json:"end_date_time"`
}

func BuildExamRoutes(router *gin.Engine) {
	router.GET("/exams/all", GetAllExams)
}

func GetAllExams(ret *gin.Context) {
	e, err := GetExamById(1)

	if err != nil {
		fmt.Println("Error getting exams!")
		fmt.Println(err)
		ret.JSON(http.StatusInternalServerError, err)
		return
	}

	ret.JSON(http.StatusOK, e)
}
func GetExamById(id int) (*exam, error) {
	s, err := GetSyllabusById(1)

	if err != nil {
		return nil, err
	}

	e := exam{
		ID:             1,
		Syllabus:       s,
		CreateDateTime: "2022-04-12T12:38:11Z",
	}
	etc, err := GetExamTagsetByExam(e)

	if err != nil {
		return nil, err
	}

	e.ExamTagset = etc
	q, err := GetQuestionsByExam(e)

	if err != nil {
		return nil, err
	}

	e.Questions = q
	return &e, nil
}
