package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

type Task struct {
	Id      int64
	Title   string
	Details string
	Active  bool
}

type database struct {
	connection *sql.DB
}

func NewDatabase() *database {
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   os.Getenv("DBHOST") + ":" + os.Getenv("DBPORT"),
		DBName: os.Getenv("DBNAME"),
	}

	// Get a database handle.
	connection, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := connection.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	log.Println("Connected!")

	db := &database{connection: connection}
	db.migrate()

	return db
}

func (db database) migrate() {}

func (db database) InsertTask(task Task) (int64, error) {
	result, err := db.connection.Exec("INSERT INTO tasks (title, details, active) VALUES (?, ?, ?)", task.Title, task.Details, task.Active)
    if err != nil {
        return 0, fmt.Errorf("InsertTask: %v", err)
    }
    id, err := result.LastInsertId()
    if err != nil {
        return 0, fmt.Errorf("InsertTask: %v", err)
    }
    return id, nil
}

func (db database) GetTasks() ([]Task, error) {
	rows, err := db.connection.Query("SELECT * FROM tasks")
	if err != nil {
		return nil, fmt.Errorf("GetTasks: %v", err)
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.Id, &task.Title, &task.Details, &task.Active); err != nil {
			return nil, fmt.Errorf("GetTasks: %v", err)

		}
		tasks = append(tasks, task)
	}
	if err:= rows.Err(); err != nil {
		return nil, fmt.Errorf("GetTasks: %v", err)
	}
	return tasks, nil
}
