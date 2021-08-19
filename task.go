package main

import "time"

type Task struct {
	Id int
	Title       string
	Description string

	TimeIntervals []TimeInterval
}

type TimeInterval struct {
	StartTime time.Time
	EndTime *time.Time
}
