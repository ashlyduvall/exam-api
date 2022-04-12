package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type question_answer struct {
	ID              int `json:"id"`
	question        question
	IsSelected      bool   `json:"is_selected"`
	IsCorrectAnswer bool   `json:"is_correct_answer"`
	Body            string `json:"body"`
}

func BuildQuestionAnswerRoutes(router *gin.Engine) {
	router.GET("/question_answers/all", GetAllQuestionAnswers)
}

func GetAllQuestionAnswers(ret *gin.Context) {
	q, err := GetQuestionById(1)

	if err != nil {
		fmt.Printf("Error getting question answers: %v\n", err)
		ret.JSON(http.StatusInternalServerError, err)
		return
	}

	a := GetQuestionAnswersByQuestion(*q)

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

func GetQuestionAnswersByQuestionAndExam(q question, e exam) []question_answer {
	return []question_answer{
		{
			ID:              1,
			question:        q,
			IsCorrectAnswer: true,
			IsSelected:      true,
			Body:            "Here is the question answer body",
		},
	}
}
