package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type tag struct {
	ID          int `json:"id"`
	Syllabus    *syllabus
	DisplayName string `json:"display_name"`
}

// --------------------
// HTTP Methods Follow
// --------------------

func BuildTagRoutes(router *gin.Engine) {
	router.GET("/tags/all", httpGetAllTags)
	router.POST("/tags/get_or_create", httpPostGetOrCreate)
}

func httpGetAllTags(ret *gin.Context) {
	s, err := GetSyllabusById("1")

	if err != nil {
		fmt.Printf("Error getting syllabus for tag: %v\n", err)
		ret.JSON(http.StatusInternalServerError, err)
		return
	}

	t, err := GetTagsBySyllabus(s)

	if err != nil {
		fmt.Printf("Error getting syllabus for tag: %v\n", err)
		ret.JSON(http.StatusInternalServerError, err)
		return
	}

	ret.JSON(http.StatusOK, t)
}

func httpPostGetOrCreate(ret *gin.Context) {
	s, err := GetSyllabusById("1")

	if err != nil {
		fmt.Printf("Error getting syllabus for tag: %v\n", err)
		ret.JSON(http.StatusInternalServerError, err)
		return
	}

	var inputJson tag
	ret.BindJSON(&inputJson)
	displayName := inputJson.DisplayName

	t, err := GetTagBySyllabusAndDisplayName(s, displayName)

	if err == sql.ErrNoRows {
		s_t, s_err := CreateNewTag(s, displayName)
		t = s_t
		err = s_err
	}

	if err != nil {
		fmt.Printf("Error getting syllabus for tag: %v\n", err)
		ret.JSON(http.StatusInternalServerError, err)
		return
	}

	ret.JSON(http.StatusOK, t)
}

// --------------------
// Raw Methods Follow
// --------------------

func GetTagsBySyllabus(s *syllabus) (*[]*tag, error) {

	ret := make([]*tag, 0)

	rows, err := DB.Query("SELECT t.ID, t.display_name FROM tags t WHERE fk_syllabus_id=?", s.ID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		t := tag{
			Syllabus: s,
		}

		err := rows.Scan(&t.ID, &t.DisplayName)

		if err != nil {
			return nil, err
		}

		ret = append(ret, &t)
	}

	return &ret, nil
}

func GetTagsByQuestion(q *question) (*[]*tag, error) {

	ret := make([]*tag, 0)

	rows, err := DB.Query(`
		SELECT t.id
		     , t.display_name 
			FROM tags t 
		 INNER JOIN question_tags qt ON qt.fk_tag_id = t.id
		 WHERE qt.fk_question_id=?
	`, q.ID)

	if err != nil {
		fmt.Printf("Error getting tags for question %v, %v", q.ID, err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		t := tag{
			Syllabus: q.Syllabus,
		}

		err := rows.Scan(&t.ID, &t.DisplayName)

		if err != nil {
			fmt.Printf("Error getting tags for question %v, %v", q.ID, err)
			return nil, err
		}

		ret = append(ret, &t)
	}

	return &ret, nil
}

func GetTagsByExam(e *exam) (*[]*tag, error) {

	ret := make([]*tag, 0)

	rows, err := DB.Query(`
		SELECT t.id
		     , t.display_name 
			FROM tags t 
		 INNER JOIN exam_tags et ON et.fk_tag_id = t.id
		 WHERE et.fk_exam_id=?
	`, e.ID)

	if err != nil {
		fmt.Printf("Error getting tags for exam %v, %v", e.ID, err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		t := tag{
			Syllabus: e.Syllabus,
		}

		err := rows.Scan(&t.ID, &t.DisplayName)

		if err != nil {
			fmt.Printf("Error getting tags for exam %v, %v", e.ID, err)
			return nil, err
		}

		ret = append(ret, &t)
	}

	return &ret, nil
}

func GetTagBySyllabusAndDisplayName(s *syllabus, d string) (*tag, error) {

	ret := tag{
		Syllabus: s,
	}

	row := DB.QueryRow(`
		SELECT t.id
		     , t.display_name 
			FROM tags t 
		 WHERE t.fk_syllabus_id = ?
		   AND t.display_name = ?
		 LIMIT 1
	`, s.ID, d)

	err := row.Scan(&ret.ID, &ret.DisplayName)

	if err != nil {
		return nil, err
	}

	return &ret, nil
}

func CreateNewTag(s *syllabus, d string) (*tag, error) {
	ret := tag{}
	ret.Syllabus = s
	ret.DisplayName = d
	fmt.Printf("Creating new tag for Syllabus %v, %v", s.ID, d)

	tx, err := DB.Begin()

	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	result, err := tx.Exec("INSERT INTO tags (fk_syllabus_id, display_name) VALUES (?, ?)", s.ID, d)

	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	ret.ID = int(id)

	return &ret, nil
}
