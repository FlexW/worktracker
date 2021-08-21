package main

import "time"

type Task struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`

	TimeIntervals []TimeInterval `json:"time_intervals"`
}

type TimeInterval struct {
	StartTime time.Time  `json:"start_time"`
	EndTime   *time.Time `json:"end_time"`
}
