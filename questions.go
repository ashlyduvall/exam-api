package main

import "github.com/gin-gonic/gin"
import "net/http"

type question struct {
	ID              int               `json:"id"`
	Syllabus        syllabus          `json:"syllabus"`
	Body            string            `json:"body"`
	QuestionAnswers []question_answer `json:"question_answers"`
}

func BuildQuestionRoutes(router *gin.Engine) {
	router.GET("/questions/get/all", GetAllQuestions)
}

func GetAllQuestions(ret *gin.Context) {
	q := GetQuestionById(1)

	ret.JSON(http.StatusOK, q)
}

func GetQuestionById(id int) question {
	s := GetSyllabusById(1)
	q := question{
		ID:       1,
		Syllabus: s,
		Body:     "Here's some question text",
	}
	q.QuestionAnswers = GetQuestionAnswersByQuestion(q)

	return q
}
