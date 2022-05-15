USE exam;

-- Syllabus
INSERT IGNORE INTO syllabus VALUES (1, 'Test Syllabus');

-- Tags
INSERT IGNORE INTO tags VALUES
  (1, 1, "Tag 1"),
  (2, 1, "Tag 2"),
  (3, 1, "Tag 3"),
  (4, 1, "Tag 4"),
  (5, 1, "Tag 5")
;

-- Questions
INSERT IGNORE INTO questions VALUES
  (1, 1, "Here is the body of question 1"),
  (2, 1, "Here is the body of question 2"),
  (3, 1, "Here is the body of question 3"),
  (4, 1, "Here is the body of question 4"),
  (5, 1, "Here is the body of question 5")
;

-- Question tags
INSERT IGNORE INTO question_tags VALUES
  (1, 1),
  (2, 2),
  (3, 3),
  (4, 4),
  (5, 5)
;

-- Question answers
INSERT IGNORE INTO question_answers (id, fk_question_id, is_correct_answer, body) VALUES
  (1, 1, true, "Correct answer to question 1"), (2, 1, false, "Incorrect answer to question 1"),
  (3, 2, true, "Correct answer to question 2"), (4, 2, false, "Incorrect answer to question 2"),
  (5, 3, true, "Correct answer to question 3"), (6, 3, false, "Incorrect answer to question 3"),
  (7, 4, true, "Correct answer to question 4"), (8, 4, false, "Incorrect answer to question 4"),
  (9, 5, true, "Multichoice correct answer to question 5"), (10, 5, false, "Incorrect answer to question 5"), (11, 5, true, "Multichoice correct answer to question 5")
;

-- Exams
INSERT IGNORE INTO exams VALUES
  (1, 1, NOW(), NULL, NULL),
  (2, 1, NOW(), NOW(), NULL),
  (3, 1, NOW(), NOW(), NOW())
;

-- Exam Tagset Tags
INSERT IGNORE INTO exam_tags VALUES
  (1, 1),
  (2, 3), (2, 4),
  (3, 5)
;

-- Exam Questions
INSERT IGNORE INTO exam_questions VALUES
  (1, 1),
  (2, 3), (2, 4),
  (3, 5)
;

-- Exam Question Answers
INSERT IGNORE INTO exam_question_answers VALUES
  (2, 3, 5), (2, 4, 8), (3, 5, 9), (3, 5, 10)
;
