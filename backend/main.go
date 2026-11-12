package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

type Quiz struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type Question struct {
	ID       int      `json:"id"`
	Question string   `json:"question"`
	Options  []string `json:"options"`
}

type Submission struct {
	QuizID      int         `json:"quizId"`
	StudentName string      `json:"studentName"`
	Answers     map[int]int `json:"answers"`
}

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == ""  {
		dbURL = "postgres://quizuser:quizpass@localhost:5432/quizdb"
	}

	mux := http.NewServeMux()

	connect := func (ctx context.Context) (*pgx.Conn, error)  {
		return pgx.Connect(ctx, dbURL)
	}

	mux.HandleFunc("GET /api/quizzes", func(w http.ResponseWriter, r *http.Request) {
		db, err := connect(r.Context())
		if err != nil { http.Error(w, "Database unavailable", 500); return }
		defer db.Close(r.Context())

		rows, err := db.Query(r.Context(), "SELECT id,  title FROM quizzes ORDER BY id")
		if err != nil { http.Error(w, "Query failed",  500); return }
		defer rows.Close()

		quizzes := []Quiz{}
		for rows.Next() {
			var q Quiz
			if err := rows.Scan(&q.ID, &q.Title); err != nil {
				http.Error(w, "Read failed", 500); return
			}
			quizzes = append(quizzes, q)
		}
		if rows.Err() != nil { http.Error(w, "Read failed",  500); return }
		writeJSON(w, quizzes)
	})

	mux.HandleFunc("GET /api/quizzes/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil { http.Error(w, "Invalid quiz ID",  400); return }

		db, err := connect(r.Context())
		if err != nil { http.Error(w, "Database unavailable",  500); return }
		defer db.Close(r.Context())

		rows, err := db.Query(r.Context(),
			"SELECT id, question, options FROM questions WHERE quiz_id=$1 ORDER BY id", id)
		if err != nil { http.Error(w, "Query failed", 500); return }
		defer rows.Close()

		questions := []Question{}
		for rows.Next() {
			var q Question
			if err := rows.Scan(&q.ID,  &q.Question, &q.Options); err != nil {
				http.Error(w, "Read failed",  500); return
			}
			questions = append(questions, q)
			if rows.Err() != nil { http.Error(w, "Read failed", 500); return }
			if len(questions) == 0 { http.Error(w, "Quiz not found", 404); return }
			writeJSON(w, questions)
		}
	})

	mux.HandleFunc("POST /api/submit", func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 64<<10)
		var s Submission
		if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
			http.Error(w, "Invalid JSON", 400); return
		}
		s.StudentName = strings.TrimSpace(s.StudentName)
		if s.QuizID <= 0 || len(s.StudentName) > 80 {
			http.Error(w, "Invalid submission", 400); return
		}

		db, err := connect(r.Context())
		if err != nil { http.Error(w,  "Database unavailable",  500); return }
		defer db.Close(r.Context())

		rows, err := db.Query(r.Context(),
			"SELECT id, correct_answer FROM questions WHERE quiz_id=$1", s.QuizID)
		if err != nil  { http.Error(w, "Query failed", 500); return }

		score, total := 0, 0
		for rows.Next() {
			var id, correct int
			if err := rows.Scan(&id, &correct); err != nil {
				rows.Close()
				http.Error(w, "Read failed", 500); return
			}
			total++
			if answer, ok := s.Answers[id]; ok && answer == correct {
				score++
			}
		}
		err = rows.Err()
		rows.Close()
		if err != nil { http.Error(w, "Read failed", 500); return }
		if total == 0 { http.Error(w, "Quiz not found", 404); return }

		_, err = db.Exec(r.Context(),
			`INSERT INTO attempts (quiz_id, student_name, score, total)
			VALUES ($1,$2,$3,$4)`,
		s.QuizID, s.StudentName, score, total)
		if err != nil { http.Error(w, "Save failed", 500); return }

		writeJSON(w, map[string]any{
			"score": score, "total": total,
		})
	})

	mux.HandleFunc("GET /api/analytics", func(w http.ResponseWriter, r *http.Request) {
		db, err := connect(r.Context())
		if err != nil { http.Error(w, "Database unavailable",  500); return }
		defer db.Close(r.Context())

		rows, err := db.Query(r.Context(),
			`SELECT q.title, ROUND(AVG(100.0*a.score/a.total),1)
			FROM attempts a JOIN quizzes q ON a.quiz_id=q.id
			WHERE a.total > 0 GROUP BY q.id, q.title ORDER BY q.id`)
		if err !=  nil { http.Error(w, "Query failed", 500); return }
		defer rows.Close()

		data := []map[string]any{}
		for rows.Next() {
			var title string
			var average float64
			if err := rows.Scan(&title, &average); err != nil {
				http.Error(w, "Read failed", 500); return
			}
			data = append(data, map[string]any{
				"title": title, "average": average,
			})
		}
		if rows.Err() != nil { http.Error(w, "Read failed", 500); return }
		writeJSON(w, data)
	})

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := os.Getenv("FRONTEND_ORIGIN")
		if origin == "" { origin = "http://localhost:5173" }
		if r.Header.Get("Origin") == origin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		mux.ServeHTTP(w,r)
	})

	server := &http.Server{
		Addr: ":8080",  Handler: handler,
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Println("API runing on port 8080")
	log.Fatal(server.ListenAndServe())
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}