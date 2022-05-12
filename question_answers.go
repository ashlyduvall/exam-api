package main

type question_answer struct {
	ID              int `json:"id"`
	question        *question
	IsSelected      bool   `json:"is_selected"`
	IsCorrectAnswer bool   `json:"is_correct_answer"`
	Body            string `json:"body"`
}

func GetQuestionAnswersByQuestion(q *question) (*[]*question_answer, error) {
	question_answers := make([]*question_answer, 0)

	rows, err := DB.Query(`
		SELECT id
		     , is_correct_answer
				 , body 
			FROM question_answers 
		 WHERE fk_question_id=?
		   AND is_deleted IS FALSE
	`, q.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		qa := question_answer{
			question: q,
		}

		err := rows.Scan(&qa.ID, &qa.IsCorrectAnswer, &qa.Body)
		if err != nil {
			return nil, err
		}
		question_answers = append(question_answers, &qa)
	}

	return &question_answers, nil
}

func GetQuestionAnswersByQuestionAndExam(q *question, e *exam) (*[]*question_answer, error) {
	question_answers := make([]*question_answer, 0)

	rows, err := DB.Query(`
		SELECT qa.id
		     , qa.is_correct_answer
				 , qa.body 
				 , eqa.fk_selected_answer_id IS NOT NULL AS is_selected
			FROM question_answers qa 
			LEFT JOIN exam_question_answers eqa
			  ON eqa.fk_question_id = qa.fk_question_id
			 AND eqa.fk_exam_id = ?
			 AND eqa.fk_selected_answer_id = qa.id
		 WHERE qa.fk_question_id=?
	`, e.ID, q.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		qa := question_answer{
			question: q,
		}

		err := rows.Scan(&qa.ID, &qa.IsCorrectAnswer, &qa.Body, &qa.IsSelected)
		if err != nil {
			return nil, err
		}
		question_answers = append(question_answers, &qa)
	}

	return &question_answers, nil
}
