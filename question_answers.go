package main

import "github.com/gin-gonic/gin"
import "net/http"

type question_answer struct {
	ID              int `json:"id"`
	question        question
	IsCorrectAnswer bool   `json:"is_correct_answer"`
	Body            string `json:"body"`
}

func BuildQuestionAnswerRoutes(router *gin.Engine) {
	router.GET("/question_answers/get/question_id/", GetAllQuestionAnswers)
}

func GetAllQuestionAnswers(ret *gin.Context) {
	q := GetQuestionById(1)
	a := GetQuestionAnswersByQuestion(q)

	ret.JSON(http.StatusOK, a)
}

func GetQuestionAnswersByQuestion(q question) []question_answer {
	return []question_answer{
		{
			ID:              1,
			question:        q,
			IsCorrectAnswer: true,
			Body:            "Here is the question answer body",
		},
	}
}
