package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type exam struct {
	ID               int          `json:"id"`
	Syllabus         *syllabus    `json:"syllabus"`
	Tags             *[]*tag      `json:"tags"`
	Questions        *[]*question `json:"questions"`
	CreateDateTime   string       `json:"create_date_time"`
	StartDateTime    string       `json:"start_date_time"`
	CompleteDateTime string       `json:"complete_date_time"`
}

func BuildExamRoutes(router *gin.Engine) {
	router.GET("/exams/all/in_progress", httpGetAllExamsInProgress)
	router.GET("/exams/all/finished", httpGetAllExamsFinished)
	router.GET("/exams/:id", httpGetExamById)
	router.POST("/exams/save", httpPostSaveExam)
	router.POST("/exams/finish", httpPostFinishExam)
	router.POST("/exams/new", httpPostNewExam)
}

// -------------------
// HTTP Methods Follow
// -------------------

func httpGetAllExamsInProgress(ret *gin.Context) {
	s, _ := GetSyllabusById("1")
	GetFinishedExams := false
	e, err := GetExamsBySyllabus(s, GetFinishedExams)

	if err != nil {
		fmt.Println("Error getting exams!")
		fmt.Println(err)
		ret.JSON(http.StatusInternalServerError, err)
		return
	}

	ret.JSON(http.StatusOK, e)
}

func httpGetAllExamsFinished(ret *gin.Context) {
	s, _ := GetSyllabusById("1")
	GetFinishedExams := true
	e, err := GetExamsBySyllabus(s, GetFinishedExams)

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

func httpPostFinishExam(ret *gin.Context) {
	var e exam
	ret.BindJSON(&e)

	err := SetExamFinished(&e)

	if err != nil {
		fmt.Println("Error finishing exam!")
		fmt.Println(err)
		ret.JSON(http.StatusInternalServerError, err)
		return
	}

	ret.JSON(http.StatusOK, gin.H{"message": "Exam Finished!"})
}

func httpPostNewExam(ret *gin.Context) {
	var ts []tag
	var e exam
	ret.BindJSON(&ts)

	s, err := GetSyllabusById("1")

	if err != nil {
		ret.JSON(http.StatusInternalServerError, err)
		return
	}

	e.Syllabus = s

	err = NewExam(&e)

	if err != nil {
		ret.JSON(http.StatusInternalServerError, err)
		return
	}

	err = SetExamTags(&e, &ts)

	if err != nil {
		ret.JSON(http.StatusInternalServerError, err)
		return
	}

	err = SetExamQuestions(&e)

	if err != nil {
		ret.JSON(http.StatusInternalServerError, err)
		return
	}

	exam, err := GetExamById(strconv.Itoa(e.ID))

	if err != nil {
		ret.JSON(http.StatusInternalServerError, err)
		return
	}

	ret.JSON(http.StatusOK, exam)
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
	`, id).Scan(&e.ID, &syllabus_id, &e.CreateDateTime, &e.StartDateTime, &e.CompleteDateTime)

	if err != nil {
		return nil, err
	}

	s, err := GetSyllabusById(syllabus_id)
	if err != nil {
		return nil, err
	}

	e.Syllabus = s

	etc, err := GetTagsByExam(&e)

	if err != nil {
		return nil, err
	}

	e.Tags = etc
	q, err := GetQuestionsByExam(&e)

	if err != nil {
		return nil, err
	}

	e.Questions = q
	return &e, nil
}

func GetExamsBySyllabus(s *syllabus, get_finished bool) (*[]*exam, error) {
	exams := make([]*exam, 0)

	var get_finished_sql string

	if get_finished {
		get_finished_sql = "WHERE complete_date_time IS NOT NULL"
	} else {
		get_finished_sql = "WHERE complete_date_time IS NULL"
	}

	sql := fmt.Sprintf(`
		SELECT id
				 , create_date_time
				 , IFNULL(start_date_time, '')
				 , IFNULL(complete_date_time, '')
			FROM exams 
			%v
	`, get_finished_sql)

	rows, err := DB.Query(sql)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		e := exam{
			Syllabus: s,
		}

		err := rows.Scan(&e.ID, &e.CreateDateTime, &e.StartDateTime, &e.CompleteDateTime)

		if err != nil {
			return nil, err
		}

		etc, err := GetTagsByExam(&e)

		if err != nil {
			return nil, err
		}

		e.Tags = etc
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

func SetExamFinished(e *exam) error {
	tx, err := DB.Begin()

	if err != nil {
		return err
	}

	defer tx.Rollback()

	// Remove all existing answers for this exam
	_, err = tx.Exec(`
    UPDATE exams 
		   SET complete_date_time = NOW()
		 WHERE id = ?
	`, e.ID)

	if err != nil {
		return err
	}

	return tx.Commit()
}

func NewExam(e *exam) error {
	tx, err := DB.Begin()

	if err != nil {
		return err
	}

	defer tx.Rollback()

	result, err := tx.Exec(`
    INSERT INTO exams (fk_syllabus_id, create_date_time, start_date_time)
		 VALUES (?, NOW(), NOW())
	`, e.Syllabus.ID)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return err
	}

	e.ID = int(id)

	return tx.Commit()
}

func SetExamTags(e *exam, ts *[]tag) error {
	tx, err := DB.Begin()

	if err != nil {
		return err
	}

	defer tx.Rollback()

	// Remove all existing tags for this exam
	_, err = tx.Exec(`
		DELETE FROM exam_tags WHERE fk_exam_id=?
	`, e.ID)

	if err != nil {
		return err
	}

	for _, t := range *ts {
		_, err := tx.Exec(`
      INSERT INTO exam_tags (fk_exam_id, fk_tag_id)
			  VALUES (?, ?)
		`, e.ID, t.ID)

		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func SetExamQuestions(e *exam) error {
	tx, err := DB.Begin()
	exam_question_limit := 50 // @TODO

	if err != nil {
		return err
	}

	defer tx.Rollback()

	// Remove all existing questions for this exam
	_, err = tx.Exec(`
		DELETE
		  FROM exam_question_answers
		 WHERE fk_exam_id=?
	`, e.ID)

	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		DELETE
		  FROM exam_questions
		 WHERE fk_exam_id=?
	`, e.ID)

	if err != nil {
		return err
	}

	_, err = tx.Exec(`
    INSERT INTO exam_questions (fk_exam_id, fk_question_id)
		 SELECT DISTINCT e.id
		      , qt.fk_question_id
		   FROM exams e
      INNER JOIN exam_tags et ON et.fk_exam_id = e.id
		  INNER JOIN question_tags qt ON qt.fk_tag_id = et.fk_tag_id
			WHERE e.id = ?
			ORDER BY MD5(CONCAT(e.id,'_',qt.fk_question_id))
			LIMIT ?
	`, e.ID, exam_question_limit)

	if err != nil {
		return err
	}

	return tx.Commit()
}
