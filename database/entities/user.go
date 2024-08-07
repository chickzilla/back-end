package entities

import "time"

type User struct {
	ID         int
	Email      string
	Password   string
	OnlySSO    bool
	Updated_at time.Time
	CreatedAt  time.Time
}
