package main

import (
	"html/template"
	"log"
	"net/http"
)

type Task struct {
	Title  string
	Detail string
	TimeSpent string
	Icon string
}

var tasks = []Task{
	{Title: "Fix bug #2312", Detail: "Needs to done a specific way.", TimeSpent: "1 hour", Icon: "bi-clock"},
	{Title: "Implement new feature", Detail: "Feature needs to be crazy fancy", TimeSpent: "30 minutes", Icon: "bi-check"},
	{Title: "Fix the docs", Detail: "Docs need update.", TimeSpent: "2 hours", Icon: "bi-check"},
	{Title: "Triage issues on GitHub", Detail: "", TimeSpent: "15 minutes", Icon: "bi-check"},
}

var templates = template.Must(template.ParseFiles("head.html", "tasks.html", "header.html", "footer.html", "reports.html"))

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "tasks.html", tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("static/css"))))

	http.Handle("/", http.RedirectHandler("/tasks/", 302))
	http.HandleFunc("/tasks/", tasksHandler)

	log.Fatal(http.ListenAndServe(":12345", nil))
}
