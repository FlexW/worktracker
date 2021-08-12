package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type TaskUi struct {
	Title     string
	Detail    string
	TimeSpent string
	Icon      string
}

type worktracker struct {
	Database *database
}

var db *database

var templates = template.Must(template.ParseFiles("head.html", "tasks.html", "header.html", "footer.html", "reports.html", "about.html", "new_task.html", "task_detail.html"))

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.GetTasks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var uiTasks []TaskUi
	for _, task := range tasks {
		uiTasks = append(uiTasks, TaskUi{Title: task.Title, Detail: task.Details, TimeSpent: "1 hour", Icon: "bi-check"})
	}

	err = templates.ExecuteTemplate(w, "tasks.html", struct {
		Tasks               []TaskUi
		HeaderTabActiveName string
	}{uiTasks, "tasks"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func reportsHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "reports.html", struct {
		HeaderTabActiveName string
	}{"reports"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "about.html", struct {
		HeaderTabActiveName string
	}{"about"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func saveTask(name string, details string) {
	db.InsertTask(Task{Title: name, Details: details})
	log.Println("Task inserted")
}

func newTaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := templates.ExecuteTemplate(w, "new_task.html", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	case "POST":
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		title := r.FormValue("title")
		details := r.FormValue("details")
		saveTask(title, details)
		http.Redirect(w, r, "/tasks/", 302)
	}
}

func taskDetailHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "task_detail.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	db = NewDatabase()

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("static/css"))))

	http.Handle("/", http.RedirectHandler("/tasks/", 302))
	http.HandleFunc("/tasks/", tasksHandler)
	http.HandleFunc("/new_task/", newTaskHandler)
	http.HandleFunc("/reports/", reportsHandler)
	http.HandleFunc("/about/", aboutHandler)
	http.HandleFunc("/task_detail/", taskDetailHandler)

	log.Fatal(http.ListenAndServe(":80", nil))
}
