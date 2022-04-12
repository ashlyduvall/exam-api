package main

type syllabus struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func GetSyllabusById(id int) syllabus {
	return syllabus{
		ID:   1,
		Name: "Test Syllabus",
	}
}
