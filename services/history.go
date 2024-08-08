package services

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/Her_feeling/back-end/database"
	"github.com/Her_feeling/back-end/database/entities"
	"github.com/gin-gonic/gin"
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

func GetUserHistories(context *gin.Context) {
	var DB = database.DB

	userId, _ := context.Get("userId")

	rows, err := DB.Query("SELECT * FROM user_history WHERE user_id = ?", userId.(int))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer rows.Close()

	var histories []entities.UserHistory

	for rows.Next() {
		var history entities.UserHistory
		if err := rows.Scan(
			&history.ID,
			&history.UserId,
			&history.Prompt,
			&history.LoveProb,
			&history.SadnessProb,
			&history.JoyProb,
			&history.AngryProb,
			&history.FearProb,
			&history.SurpriseProb,
			&history.CreatedAt); err != nil {

			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		histories = append(histories, history)
	}

	context.JSON(http.StatusOK, gin.H{"data": histories})

}
