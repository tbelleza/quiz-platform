CREATE TABLE quizzes (
  id SERIAL PRIMARY KEY,
  title TEXT NOT NULL
);

CREATE TABLE questions (
  id SERIAL PRIMARY KEY,
  quiz_id INTEGER REFERENCES quizzes(id),
  question TEXT NOT NULL,
  options JSONB NOT NULL,
  correct_answer INTEGER NOT NULL
);

CREATE TABLE attempts (
  id SERIAL PRIMARY KEY,
  quiz_id INTEGER REFERENCES quizzes(id),
  student_name TEXT NOT NULL,
  score INTEGER NOT NULL,
  total INTEGER NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

INSERT INTO quizzes (title) VALUES
('Go Basics'),
('Web Development');

INSERT INTO questions
(quiz_id, question, options, correct_answer)
VALUES
(1, 'Which keyword declares a variable in Go?',
 '["let", "var", "const", "define"]', 1),
(1, 'Which package is used for formatted output?',
 '["fmt", "http", "os", "math"]', 0),
(1, 'What is the file extension for Go?',
 '[".go", ".golang", ".g", ".gojs"]', 0),
(2, 'What does HTML stand for?',
 '["HyperText Markup Language", "High Tech Machine Language", "Home Tool Markup Language"]', 0),
(2, 'Which framework are we using?',
 '["React", "Angular", "Vue"]', 2),
(2, 'What language adds types to JavaScript?',
 '["Python", "TypeScript", "R"]', 1);