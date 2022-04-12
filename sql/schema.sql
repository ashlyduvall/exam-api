-- DB Schema for the exam app

DROP DATABASE IF EXISTS exam;
CREATE DATABASE IF NOT EXISTS exam;
USE exam;

CREATE TABLE IF NOT EXISTS syllabus (
  id INT(11) NOT NULL PRIMARY KEY AUTO_INCREMENT,
  display_name VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS tags (
  id INT(11) NOT NULL PRIMARY KEY AUTO_INCREMENT,
  fk_syllabus_id INT(11) NOT NULL,
  display_name VARCHAR(255) NOT NULL,
  CONSTRAINT FOREIGN KEY (fk_syllabus_id) REFERENCES syllabus (id)
);

CREATE TABLE IF NOT EXISTS questions (
  id INT(11) NOT NULL PRIMARY KEY AUTO_INCREMENT,
  fk_syllabus_id INT(11) NOT NULL,
  body TEXT,
  CONSTRAINT FOREIGN KEY (fk_syllabus_id) REFERENCES syllabus (id)
);

CREATE TABLE IF NOT EXISTS question_tags (
  fk_question_id INT(11) NOT NULL,
  fk_tag_id INT(11) NOT NULL,
  PRIMARY KEY(fk_question_id, fk_tag_id),
  CONSTRAINT FOREIGN KEY (fk_question_id) REFERENCES questions (id),
  CONSTRAINT FOREIGN KEY (fk_tag_id) REFERENCES tags (id)
);

CREATE TABLE IF NOT EXISTS question_answers (
  id INT(11) NOT NULL PRIMARY KEY AUTO_INCREMENT,
  fk_question_id INT(11) NOT NULL,
  is_correct_answer BOOLEAN DEFAULT FALSE,
  body VARCHAR(2048) NOT NULL,
  CONSTRAINT FOREIGN KEY (fk_question_id) REFERENCES questions (id)
);

CREATE TABLE IF NOT EXISTS exam_tagsets (
  id INT(11) NOT NULL PRIMARY KEY AUTO_INCREMENT,
  fk_syllabus_id INT(11) NOT NULL,
  display_name VARCHAR(255) NOT NULL,
  CONSTRAINT FOREIGN KEY (fk_syllabus_id) REFERENCES syllabus (id)
);

CREATE TABLE IF NOT EXISTS exam_tagset_tags (
  fk_exam_tagset_id INT(11) NOT NULL,
  fk_tag_id INT(11) NOT NULL,
  PRIMARY KEY(fk_exam_tagset_id, fk_tag_id),
  CONSTRAINT FOREIGN KEY (fk_exam_tagset_id) REFERENCES exam_tagsets (id),
  CONSTRAINT FOREIGN KEY (fk_tag_id) REFERENCES tags (id)
);

CREATE TABLE IF NOT EXISTS exams (
  id INT(11) NOT NULL PRIMARY KEY AUTO_INCREMENT,
  fk_syllabus_id INT(11) NOT NULL,
  fk_exam_tagset_id INT(11) NOT NULL,
  create_date_time DATETIME NOT NULL,
  start_date_time DATETIME NULL DEFAULT NULL,
  complete_date_time DATETIME NULL DEFAULT NULL,
  CONSTRAINT FOREIGN KEY (fk_syllabus_id) REFERENCES syllabus (id),
  CONSTRAINT FOREIGN KEY (fk_exam_tagset_id) REFERENCES exam_tagsets (id)
);

CREATE TABLE IF NOT EXISTS exam_questions (
  fk_exam_id INT(11) NOT NULL,
  fk_question_id INT(11) NOT NULL,
  PRIMARY KEY (fk_exam_id, fk_question_id),
  CONSTRAINT FOREIGN KEY (fk_exam_id) REFERENCES exams (id),
  CONSTRAINT FOREIGN KEY (fk_question_id) REFERENCES questions (id)
);

CREATE TABLE IF NOT EXISTS exam_question_answers (
  fk_exam_id INT(11) NOT NULL,
  fk_question_id INT(11) NOT NULL,
  fk_selected_answer_id INT(11) NOT NULL,
  PRIMARY KEY (fk_exam_id, fk_question_id, fk_selected_answer_id),
  CONSTRAINT FOREIGN KEY (fk_exam_id) REFERENCES exams (id),
  CONSTRAINT FOREIGN KEY (fk_question_id) REFERENCES questions (id),
  CONSTRAINT FOREIGN KEY (fk_selected_answer_id) REFERENCES question_answers (id)
);
