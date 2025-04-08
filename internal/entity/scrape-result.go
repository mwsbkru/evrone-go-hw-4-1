package entity

import "time"

type ScrapeResult struct {
	Date           time.Time
	Url            string
	StatusCode     int
	Title          string
	Description    string
	SuccessAttempt int // Номер попытки, с какой получилось загрузить данные
}
