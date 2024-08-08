package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"

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
type requestText struct {
	Prompt string `json:"prompt" binding:"required"`
}

func SendPrompt(c *gin.Context) (map[string]float64, error) {
	var promptRequest requestText
	if err := c.ShouldBindJSON(&promptRequest); err != nil {
		return nil, errors.New("can't bind json")
	}

	envChan := make(chan string, 1)
	errChan := make(chan error, 1)
	endCodedCh := make(chan string, 1)
	var wg sync.WaitGroup
	wg.Add(2)

	go utils.GetEnv(&wg, envChan, errChan)

	go func() {
		defer wg.Done()
		encodedPrompt := url.QueryEscape(promptRequest.Prompt)
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
		err := CreateUserHistory(emailString, promptRequest.Prompt, responseData)
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
