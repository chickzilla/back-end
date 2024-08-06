package entities

import "time"

type User struct {
	ID         int
	Email      string
	Password   string
	Updated_at time.Time
	CreatedAt  time.Time
}

type UserHistory struct {
	ID           int
	UserId       int
	Prompt       string
	SadnessProb  float64
	LoveProb     float64
	JoyProb      float64
	AngryProb    float64
	FearProb     float64
	SurpriseProb float64
	CreatedAt    time.Time
}
