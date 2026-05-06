package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

const (
	basePattern string = "/"
)

type Todo struct {
	Id   int
	Text string
	Done bool
}

var db *sql.DB

var tmpl = template.Must(template.ParseFiles("templates/index.html"))
var todos []Todo

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "todo.db")
	if err != nil {
		log.Fatal(err)
	}

	createTable := `
		CREATE TABLE IF NOT EXISTS todos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			text TEXT,
			done INTEGER
		);
	`

	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	initDB()

	http.HandleFunc(basePattern, indexHandler)
	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/delete", deleteHandler)
	http.HandleFunc("/toggle", toggleHandler)
	http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	selectQuery := `
		SELECT id, text, done FROM todos ORDER BY id DESC
	`
	rows, err := db.Query(selectQuery)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var t Todo
		var done int
		rows.Scan(&t.Id, &t.Text, &done)
		t.Done = done == 1
		todos = append(todos, t)
	}

	tmpl.Execute(w, todos)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		task := r.FormValue("task")

		insertQuery := `
			INSERT INTO todos (text, done) VALUES (?, ?)
		`
		_, err := db.Exec(insertQuery, task, false)

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}
	http.Redirect(w, r, basePattern, http.StatusSeeOther)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		id := r.FormValue("id")

		deleteQuery := `
			DELETE FROM todos WHERE id = ?
		`
		_, err := db.Exec(deleteQuery, id)
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
	}
	http.Redirect(w, r, basePattern, http.StatusSeeOther)
}

func toggleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		id := r.FormValue("toggle")
		changeQuery := `
			UPDATE todos 
				SET done = CASE WHEN done = 1 THEN 0 ELSE 1 END
				WHERE id = ?
		`
		_, err := db.Exec(changeQuery, id)

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}
	http.Redirect(w, r, basePattern, http.StatusSeeOther)
}
