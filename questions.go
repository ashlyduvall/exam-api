package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type question struct {
	ID              int               `json:"id"`
	Syllabus        *syllabus         `json:"syllabus"`
	Body            string            `json:"body"`
	Tags            *[]*tag           `json:"tags"`
	QuestionAnswers []question_answer `json:"question_answers"`
}

// --------------------
// HTTP Methods Follow
// --------------------

func BuildQuestionRoutes(router *gin.Engine) {
	router.GET("/questions/all", httpGetAllQuestions)
	router.GET("/questions/:id", httpGetQuestionById)
}

func httpGetAllQuestions(ret *gin.Context) {
	q, err := GetQuestionById("1")
	if err != nil {
		fmt.Printf("Error getting all questions: %v\n", err)
		ret.JSON(http.StatusInternalServerError, err)
		return
	}

	ret.JSON(http.StatusOK, []*question{q})
}

func httpGetQuestionById(ret *gin.Context) {
	question_id := ret.Param("id")
	q, err := GetQuestionById(question_id)
	if err != nil {
		fmt.Printf("Error getting question: %v\n", err)
		ret.JSON(http.StatusInternalServerError, err)
		return
	}

	ret.JSON(http.StatusOK, q)
}

// -------------------
// Raw Methods Follow
// -------------------

func GetQuestionById(id string) (*question, error) {
	q := question{}
	var syllabus_id string
	err := DB.QueryRow("SELECT q.id, q.fk_syllabus_id, q.body FROM questions q WHERE id=?", id).Scan(&q.ID, &syllabus_id, &q.Body)
	if err != nil {
		return nil, err
	}

	s, err := GetSyllabusById(syllabus_id)
	if err != nil {
		return nil, err
	}

	q_tags, err := GetTagsByQuestion(q)
	if err != nil {
		return nil, err
	}

	q.Syllabus = s
	q.Tags = q_tags
	q.QuestionAnswers = GetQuestionAnswersByQuestion(q)

	return &q, nil
}

func GetQuestionsByExam(e exam) (*[]*question, error) {
	q, err := GetQuestionById("1")

	if err != nil {
		return nil, err
	}

	ql := []*question{q}

	return &ql, nil
}
