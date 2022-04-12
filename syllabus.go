package main

type syllabus struct {
	ID          int    `json:"id"`
	DisplayName string `json:"display_name"`
}

func GetSyllabusById(id string) (*syllabus, error) {
	s := syllabus{}
	err := DB.QueryRow("SELECT s.id, s.display_name FROM syllabus s WHERE s.id=?", id).Scan(&s.ID, &s.DisplayName)

	if err != nil {
		return nil, err
	}

	return &s, nil
}
