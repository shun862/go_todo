package main

import (
	"html/template"
	"net/http"
	"slices"
	"strconv"
)

type Todo struct {
	Id   int
	Text string
	Done bool
}

var tmpl = template.Must(template.ParseFiles("templates/index.html"))
var todos []Todo

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/delete", deleteHandler)
	http.HandleFunc("/toggle", toggleHandler)
	http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.Execute(w, todos)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		task := r.FormValue("task")
		todo := Todo{
			Id:   len(todos) + 1,
			Text: task,
		}
		todos = append(todos, todo)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		idStr := r.FormValue("id")
		id, _ := strconv.Atoi(idStr)
		todos = slices.DeleteFunc(todos, func(t Todo) bool {
			return t.Id == id
		})
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func toggleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		idStr := r.FormValue("toggle")
		id, _ := strconv.Atoi(idStr)
		for i, todo := range todos {
			if todo.Id == id {
				todos[i].Done = !todos[i].Done
				break
			}
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
