package services

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"

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

var allowedSortColumns = map[string]bool{
	"created_at":    true,
	"love_prob":     true,
	"sadness_prob":  true,
	"joy_prob":      true,
	"angry_prob":    true,
	"fear_prob":     true,
	"surprise_prob": true,
}

var allowedOrderDirections = map[string]bool{
	"ASC":  true,
	"DESC": true,
}

func GetUserHistories(context *gin.Context) {
	var DB = database.DB

	userId, _ := context.Get("userId")

	limitStr := context.Query("limit")
	offsetStr := context.Query("offset")
	sortByStr := context.Query("sortBy")
	orderByStr := context.Query("orderBy")

	// default
	limit := 5
	offset := 0
	sortBy := "created_at"
	orderBy := "DESC"

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			offset = o
		}
	}

	if sortByStr != "" {
		if allowedSortColumns[sortByStr] {
			sortBy = sortByStr
		}
	}
	if orderByStr != "" {
		if allowedOrderDirections[orderByStr] {
			orderBy = orderByStr
		}
	}

	fmt.Println("sortByString", sortByStr)
	fmt.Println("orderByString", orderByStr)
	fmt.Println("sortBy", sortBy)
	fmt.Println("orderBy", orderBy)

	rows, err := DB.Query(`
		SELECT * FROM user_history
		WHERE user_id = ?
		ORDER BY `+sortBy+` `+orderBy+`
		LIMIT ? OFFSET ?`, userId.(int), limit, offset)

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
		history.LoveProb = math.Round(history.LoveProb*100) / 100
		history.SadnessProb = math.Round(history.SadnessProb*100) / 100
		history.JoyProb = math.Round(history.JoyProb*100) / 100
		history.AngryProb = math.Round(history.AngryProb*100) / 100
		history.FearProb = math.Round(history.FearProb*100) / 100
		history.SurpriseProb = math.Round(history.SurpriseProb*100) / 100

		histories = append(histories, history)
	}

	var totalRecords int
	if err = DB.QueryRow("SELECT COUNT(*) FROM user_history WHERE user_id = ?", userId.(int)).Scan(&totalRecords); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	countRecords := len(histories)

	context.JSON(http.StatusOK, gin.H{"data": gin.H{"items": histories, "metaData": gin.H{"total": totalRecords, "count": countRecords}}})

}
