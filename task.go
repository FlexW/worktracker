package main

import "time"

type Task struct {
	Id          int           `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Duration    time.Duration `json:"duration"`
	Active      bool          `json:"active"`
}

type NewTask struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	StartTime   time.Time  `json:"startTime"`
	EndTime     *time.Time `json:"endTime"`
}

type TimeInterval struct {
	StartTime time.Time  `json:"startTime"`
	EndTime   *time.Time `json:"endTime"`
}
