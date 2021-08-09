package main

import (
	"html/template"
	"log"
	"net/http"
)

type Task struct {
	Title     string
	Detail    string
	TimeSpent string
	Icon      string
}

var tasks = []Task{
	{Title: "Fix bug #2312", Detail: "Needs to be done a specific way.", TimeSpent: "1 hour", Icon: "bi-clock"},
	{Title: "Implement new feature", Detail: "Feature needs to be crazy fancy", TimeSpent: "30 minutes", Icon: "bi-check"},
	{Title: "Fix the docs", Detail: "Docs need update.", TimeSpent: "2 hours", Icon: "bi-check"},
	{Title: "Triage issues on GitHub", Detail: "", TimeSpent: "15 minutes", Icon: "bi-check"},
}

var templates = template.Must(template.ParseFiles("head.html", "tasks.html", "header.html", "footer.html", "reports.html", "about.html", "new_task.html", "task_detail.html"))

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "tasks.html", struct {
		Tasks []Task
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

func newTaskHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "new_task.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func taskDetailHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "task_detail.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("static/css"))))

	http.Handle("/", http.RedirectHandler("/tasks/", 302))
	http.HandleFunc("/tasks/", tasksHandler)
	http.HandleFunc("/new_task/", newTaskHandler)
	http.HandleFunc("/reports/", reportsHandler)
	http.HandleFunc("/about/", aboutHandler)
	http.HandleFunc("/task_detail/", taskDetailHandler)

	log.Fatal(http.ListenAndServe(":12345", nil))
}
