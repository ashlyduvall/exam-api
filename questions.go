package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type question struct {
	ID              int                 `json:"id"`
	Syllabus        *syllabus           `json:"syllabus"`
	Body            string              `json:"body"`
	Tags            *[]*tag             `json:"tags"`
	QuestionAnswers *[]*question_answer `json:"question_answers"`
}

// --------------------
// HTTP Methods Follow
// --------------------

func BuildQuestionRoutes(router *gin.Engine) {
	router.GET("/questions/all", httpGetAllQuestions)
	router.GET("/questions/:id", httpGetQuestionById)
	router.POST("/questions/save", httpPostSaveQuestion)
}

func httpGetAllQuestions(ret *gin.Context) {
	s, err := GetSyllabusById("1")
	if err != nil {
		fmt.Printf("Error getting all questions - couldn't find Syllabus: %v\n", err)
		ret.JSON(http.StatusInternalServerError, err)
		return
	}

	q, err := GetQuestionsBySyllabus(s)
	if err != nil {
		fmt.Printf("Error getting all questions: %v\n", err)
		ret.JSON(http.StatusInternalServerError, err)
		return
	}

	ret.JSON(http.StatusOK, q)
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

func httpPostSaveQuestion(ret *gin.Context) {
	var q question
	ret.BindJSON(&q)
	fmt.Print(q)

	// Save tags
	SetQuestionTags(&q)

	// Save answers
	// SetQuestionAnswers(&q)

	// Save question body

	ret.JSON(http.StatusOK, q)
}

// -------------------
// Raw Methods Follow
// -------------------

func GetQuestionsBySyllabus(s *syllabus) (*[]*question, error) {
	var ret []*question

	rows, err := DB.Query("SELECT q.id, q.body FROM questions q WHERE q.fk_syllabus_id=?", s.ID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		q := question{
			Syllabus: s,
		}

		err := rows.Scan(&q.ID, &q.Body)
		if err != nil {
			return nil, err
		}

		q_tags, err := GetTagsByQuestion(&q)
		if err != nil {
			return nil, err
		}
		q.Tags = q_tags

		q_answers, err := GetQuestionAnswersByQuestion(&q)
		if err != nil {
			return nil, err
		}
		q.QuestionAnswers = q_answers

		ret = append(ret, &q)
	}

	return &ret, nil
}

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

	q_tags, err := GetTagsByQuestion(&q)
	if err != nil {
		return nil, err
	}

	q.Syllabus = s
	q.Tags = q_tags
	q.QuestionAnswers, err = GetQuestionAnswersByQuestion(&q)

	if err != nil {
		return nil, err
	}

	return &q, nil
}

func GetQuestionByIdWithExamData(id string, e *exam) (*question, error) {
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

	q_tags, err := GetTagsByQuestion(&q)
	if err != nil {
		return nil, err
	}

	q.Syllabus = s
	q.Tags = q_tags
	question_answers, err := GetQuestionAnswersByQuestionAndExam(&q, e)
	if err != nil {
		return nil, err
	}

	q.QuestionAnswers = question_answers

	return &q, nil
}

func GetQuestionsByExam(e *exam) (*[]*question, error) {
	rows, err := DB.Query(`
    SELECT eq.fk_question_id
		  FROM exam_questions eq
		 WHERE eq.fk_exam_id=?
	`, e.ID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []*question
	for rows.Next() {
		var question_id string
		err := rows.Scan(&question_id)

		if err != nil {
			return nil, err
		}

		question, err := GetQuestionByIdWithExamData(question_id, e)

		if err != nil {
			return nil, err
		}

		questions = append(questions, question)
	}

	return &questions, nil
}

func SetQuestionTags(q *question) error {
	tx, err := DB.Begin()

	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM question_tags WHERE fk_question_id=?", q.ID)

	if err != nil {
		return err
	}

	for _, t := range *q.Tags {
		_, err := tx.Exec("INSERT INTO question_tags (fk_question_id, fk_tag_id) VALUES (?, ?)", q.ID, t.ID)

		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
