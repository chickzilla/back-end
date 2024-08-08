package services

import (
	"database/sql"
	"errors"

	"github.com/Her_feeling/back-end/database"
)

func CreateUserHistory(email, prompt string, response TextResponseData) error {
	var DB = database.DB
	var userId string
	moods := response.Data

	err := DB.QueryRow("SELECT id FROM user WHERE email = ?", email).Scan(&userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("email not found")
		} else {
			return err
		}
	}

	_, err = DB.Exec("INSERT INTO user_history (user_id, prompt, love_prob, sadness_prob, joy_prob, angry_prob, surprise_prob, fear_prob) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", userId, prompt, moods.Love, moods.Sadness, moods.Joy, moods.Anger, moods.Surprise, moods.Fear)

	if err != nil {
		return err
	}

	return nil
}
