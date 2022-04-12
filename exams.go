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
	router.GET("/exams/:id", httpGetExamById)
}

// -------------------
// HTTP Methods Follow
// -------------------

func httpGetExamById(ret *gin.Context) {
	exam_id := ret.Param("id")
	e, err := GetExamById(exam_id)

	if err != nil {
		fmt.Println("Error getting exams!")
		fmt.Println(err)
		ret.JSON(http.StatusInternalServerError, err)
		return
	}

	ret.JSON(http.StatusOK, e)
}

// -------------------
// Raw Methods Follow
// -------------------

func GetExamById(id string) (*exam, error) {
	e := exam{}
	var syllabus_id string
	err := DB.QueryRow("SELECT id, fk_syllabus_id, create_date_time, IFNULL(start_date_time, ''), IFNULL(complete_date_time, '') FROM exams WHERE id=?", id).Scan(&e.ID, &syllabus_id, &e.CreateDateTime, &e.StartDateTime, &e.EndDateTime)

	if err != nil {
		return nil, err
	}

	s, err := GetSyllabusById(syllabus_id)
	if err != nil {
		return nil, err
	}

	e.Syllabus = s

	etc, err := GetExamTagsetByExam(e)

	if err != nil {
		return nil, err
	}

	e.ExamTagset = etc
	q, err := GetQuestionsByExam(&e)

	if err != nil {
		return nil, err
	}

	e.Questions = q
	return &e, nil
}
