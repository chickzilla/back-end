package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"

	"github.com/Her_feeling/back-end/database"
	utils "github.com/Her_feeling/back-end/utils/helper"
	"github.com/gin-gonic/gin"
)

type TextResponseData struct {
	Data struct {
		Sadness  float64 `json:"sadness"`
		Joy      float64 `json:"joy"`
		Love     float64 `json:"love"`
		Anger    float64 `json:"anger"`
		Fear     float64 `json:"fear"`
		Surprise float64 `json:"surprise"`
	} `json:"data"`
}

// sadness (0), joy (1), love (2), anger (3), fear (4), and surprise (5).

func SendPrompt(c *gin.Context) (map[string]float64, error) {
	prompt := c.Query("prompt")

	envChan := make(chan string, 1)
	errChan := make(chan error, 1)
	endCodedCh := make(chan string, 1)
	var wg sync.WaitGroup
	wg.Add(2)

	go utils.GetEnv(&wg, envChan, errChan)

	go func() {
		defer wg.Done()
		encodedPrompt := url.QueryEscape(prompt)
		endCodedCh <- encodedPrompt
	}()

	wg.Wait()
	close(errChan)
	close(envChan)
	close(endCodedCh)

	if err := <-errChan; err != nil {
		return nil, err
	}

	textAIURL := <-envChan

	encodedPrompt := <-endCodedCh

	usedURL := fmt.Sprintf("%s/prompt?query=%s", textAIURL, encodedPrompt)

	response, err := http.Get(usedURL)
	if err != nil {
		return nil, errors.New("can't send request to text-ai service")
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.New("can't read response body")
	}

	var responseData TextResponseData
	if err := json.Unmarshal(body, &responseData); err != nil {
		return nil, errors.New("can't unmarshal response body")
	}

	if email, ok := c.Get("email"); ok {
		emailString := email.(string)
		err := CreateUserHistory(emailString, prompt, responseData)
		fmt.Println("error when create user history: ", err)
	}

	emotions := map[string]float64{
		"sadness":  responseData.Data.Sadness,
		"joy":      responseData.Data.Joy,
		"love":     responseData.Data.Love,
		"anger":    responseData.Data.Anger,
		"fear":     responseData.Data.Fear,
		"surprise": responseData.Data.Surprise,
	}
	return emotions, nil
}

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
