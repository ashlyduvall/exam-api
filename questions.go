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

func BuildQuestionRoutes(router *gin.Engine) {
	router.GET("/questions/all", GetAllQuestions)
}

func GetAllQuestions(ret *gin.Context) {
	q, err := GetQuestionById(1)
	if err != nil {
		fmt.Printf("Error getting all questions: %v\n", err)
		ret.JSON(http.StatusInternalServerError, err)
		return
	}

	ret.JSON(http.StatusOK, q)
}

func GetQuestionById(id int) (*question, error) {
	s, err := GetSyllabusById(1)

	if err != nil {
		return nil, err
	}

	q := question{
		ID:       1,
		Syllabus: s,
		Body:     "Here's some question text",
	}
	q_tags, err := GetTagsByQuestion(q)
	if err != nil {
		return nil, err
	}
	q.Tags = q_tags
	q.QuestionAnswers = GetQuestionAnswersByQuestion(q)

	return &q, nil
}

func GetQuestionByIdAndExam(id int, e exam) (*question, error) {
	s, err := GetSyllabusById(1)

	if err != nil {
		return nil, err
	}

	q := question{
		ID:       1,
		Syllabus: s,
		Body:     "Here's some question text",
	}
	q_tags, err := GetTagsByQuestion(q)

	if err != nil {
		return nil, err
	}

	q.Tags = q_tags
	q.QuestionAnswers = GetQuestionAnswersByQuestionAndExam(q, e)

	return &q, nil
}

func GetQuestionsByExam(e exam) (*[]*question, error) {
	q, err := GetQuestionByIdAndExam(1, e)

	if err != nil {
		return nil, err
	}

	ql := []*question{q}

	return &ql, nil
}
