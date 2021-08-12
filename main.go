package main

import (
	"html/template"
	"log"
	"net/http"
	"fmt"
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

var tasks = []TaskUi{
	{Title: "Fix bug #2312", Detail: "Needs to be done a specific way.", TimeSpent: "1 hour", Icon: "bi-clock"},
	{Title: "Implement new feature", Detail: "Feature needs to be crazy fancy", TimeSpent: "30 minutes", Icon: "bi-check"},
	{Title: "Fix the docs", Detail: "Docs need update.", TimeSpent: "2 hours", Icon: "bi-check"},
	{Title: "Triage issues on GitHub", Detail: "", TimeSpent: "15 minutes", Icon: "bi-check"},
}

var db *database

var templates = template.Must(template.ParseFiles("head.html", "tasks.html", "header.html", "footer.html", "reports.html", "about.html", "new_task.html", "task_detail.html"))

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "tasks.html", struct {
		Tasks []TaskUi
		HeaderTabActiveName string
	}{tasks, "tasks"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		// fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", r.PostForm)
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
