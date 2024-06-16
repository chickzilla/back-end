package entities

import "time"

type User struct {
	ID        string
	Username  string
	Password  string
	Updated_at time.Time
	CreatedAt time.Time
}

