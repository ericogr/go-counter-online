package counter

import "time"

type Counter struct {
	UUID  string    `json:"uuid"`
	Name  string    `json:"name"`
	Count int       `json:"count"`
	Date  time.Time `json:"created_at"`
}
