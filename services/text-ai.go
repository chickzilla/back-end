package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type TextResponseData struct {
	Data []float64 `json:"data"`
}

// sadness (0), joy (1), love (2), anger (3), fear (4), and surprise (5).

func SendPrompt(prompt string) (map[string]float64, error) {
	envChan := make(chan string, 1)
	errChan := make(chan error, 1)
	endCodedCh := make(chan string, 1)
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		err := godotenv.Load()
		if err != nil {
			errChan <- errors.New("error loading .env file")
			return
		}

		textAIURL, isInEnv := os.LookupEnv("AI_TEXT_URL")
		if !isInEnv {
			errChan <- errors.New(".env file has no variable AI_TEXT_URL")
			
		}
		if textAIURL == "" {
			errChan <- errors.New(".env variable AI_TEXT_URL is empty")
			return
		}

		envChan <- textAIURL
	}()

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
	fmt.Println("usedURL ",usedURL )

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

	emotions := listToMapJSON(responseData.Data)
	return emotions, nil
}

func listToMapJSON(data []float64) map[string]float64 {
	emotions := map[string]float64{
		"sadness":  data[0],
		"joy":      data[1],
		"love":     data[2],
		"anger":    data[3],
		"fear":     data[4],
		"surprise": data[5],
	}
	return emotions
}
