package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
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

	// Save tags
	err := SetQuestionTags(&q)

	if err != nil {
		ret.JSON(http.StatusInternalServerError, err)
		return
	}

	// Save answers
	err = SetQuestionAnswers(&q)

	if err != nil {
		ret.JSON(http.StatusInternalServerError, err)
		return
	}

	// Save question body
	err = SetQuestionBody(&q)

	if err != nil {
		ret.JSON(http.StatusInternalServerError, err)
		return
	}

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

func SetQuestionAnswers(q *question) error {
	tx, err := DB.Begin()

	if err != nil {
		return err
	}

	defer tx.Rollback()

	var answers_to_keep []interface{}
	answers_to_keep = append(answers_to_keep, q.ID)

	for _, a := range *q.QuestionAnswers {
		var err error
		if a.ID > 0 {
			_, e := tx.Exec(`
				UPDATE question_answers
				   SET is_correct_answer=?
					   , body=?
				 WHERE id=?
			`, a.IsCorrectAnswer, a.Body, a.ID)
			err = e
			answers_to_keep = append(answers_to_keep, a.ID)
		} else {
			result, e := tx.Exec(`
        INSERT INTO question_answers (fk_question_id, is_correct_answer, body) VALUES (?, ?, ?)
			`, q.ID, a.IsCorrectAnswer, a.Body)
			err = e
			id, _ := result.LastInsertId()
			answers_to_keep = append(answers_to_keep, int(id))
		}

		if err != nil {
			return err
		}

	}

	fmt.Println(answers_to_keep)

	// Handle removing answers

	var sql string
	if len(answers_to_keep) == 0 {
		sql = `
			UPDATE question_answers
				 SET is_deleted = TRUE
			 WHERE fk_question_id = ?
    `
	} else {
		q_marks := strings.Repeat("?,", len(answers_to_keep)-1)
		q_marks = q_marks[:len(q_marks)-1]
		sql = fmt.Sprintf(`
			UPDATE question_answers
				 SET is_deleted = TRUE
			 WHERE fk_question_id = ?
		     AND id NOT IN (%v)
    `, q_marks)
	}

	_, err = tx.Exec(sql, answers_to_keep...)

	if err != nil {
		return err
	}

	return tx.Commit()
}

func SetQuestionBody(q *question) error {
	tx, err := DB.Begin()

	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.Exec(`
		UPDATE questions
			 SET body=?
		 WHERE id=?
	`, q.Body, q.ID)

	if err != nil {
		return err
	}

	return tx.Commit()
}
