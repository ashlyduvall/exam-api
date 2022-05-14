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
	router.GET("/exams/all", httpGetAllExams)
	router.GET("/exams/:id", httpGetExamById)
	router.POST("/exams/save", httpPostSaveExam)
}

// -------------------
// HTTP Methods Follow
// -------------------

func httpGetAllExams(ret *gin.Context) {
	s, _ := GetSyllabusById("1")
	e, err := GetExamsBySyllabus(s)

	if err != nil {
		fmt.Println("Error getting exams!")
		fmt.Println(err)
		ret.JSON(http.StatusInternalServerError, err)
		return
	}

	ret.JSON(http.StatusOK, e)
}

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

func httpPostSaveExam(ret *gin.Context) {
	var e exam
	ret.BindJSON(&e)

	err := SetExamQuestionAnswers(&e)

	if err != nil {
		fmt.Println("Error saving exam!")
		fmt.Println(err)
		ret.JSON(http.StatusInternalServerError, err)
		return
	}

	ret.JSON(http.StatusOK, gin.H{"message": "Exam Saved!"})
}

// -------------------
// Raw Methods Follow
// -------------------

func GetExamById(id string) (*exam, error) {
	e := exam{}
	var syllabus_id string
	err := DB.QueryRow(`
		SELECT id
				 , fk_syllabus_id
				 , create_date_time
				 , IFNULL(start_date_time, '')
				 , IFNULL(complete_date_time, '')
			FROM exams 
		 WHERE id=?
	`, id).Scan(&e.ID, &syllabus_id, &e.CreateDateTime, &e.StartDateTime, &e.EndDateTime)

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

func GetExamsBySyllabus(s *syllabus) (*[]*exam, error) {
	exams := make([]*exam, 0)

	rows, err := DB.Query(`
		SELECT id
				 , create_date_time
				 , IFNULL(start_date_time, '')
				 , IFNULL(complete_date_time, '')
			FROM exams 
	`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		e := exam{
			Syllabus: s,
		}

		err := rows.Scan(&e.ID, &e.CreateDateTime, &e.StartDateTime, &e.EndDateTime)

		if err != nil {
			return nil, err
		}

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
		exams = append(exams, &e)
	}

	return &exams, nil
}

func SetExamQuestionAnswers(e *exam) error {
	tx, err := DB.Begin()

	if err != nil {
		return err
	}

	defer tx.Rollback()

	// Remove all existing answers for this exam
	_, err = tx.Exec(`
    DELETE FROM exam_question_answers
		 WHERE fk_exam_id = ?
	`, e.ID)

	if err != nil {
		return err
	}

	// Insert current state
	for _, q := range *e.Questions {
		for _, a := range *q.QuestionAnswers {
			if a.IsSelected {
				_, err := tx.Exec(`
					INSERT INTO exam_question_answers (fk_exam_id, fk_question_id, fk_selected_answer_id)
					 VALUES (?, ?, ?)
				`, e.ID, q.ID, a.ID)
				if err != nil {
					return err
				}
			}
		}
	}

	return tx.Commit()
}
