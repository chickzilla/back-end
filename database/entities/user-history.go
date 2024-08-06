package entities

import "time"

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
