package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

type TextResponseData struct {
	Data []float64 `json:"data"`
}
// sadness (0), joy (1), love (2), anger (3), fear (4), and surprise (5).

func SendPrompt(prompt string) (map[string]float64, error) {
	
	err := godotenv.Load()
	if err != nil {
		return nil , errors.New("error loading .env file")
	}
	
	textAIURL := os.Getenv("AI_TEXT_URL")
	//fmt.Print("textAIURL : ", textAIURL)

	if textAIURL == "" {
		return nil, errors.New(".env variable AI_TEXT_URL is not set")
	}

	// ทำ text ให้อยู่ในรูปแบบ url ( ตัดพวก space --> %20)
	encodedPrompt := url.QueryEscape(prompt)
	usedURL := fmt.Sprintf("%s/prompt?query=%s", textAIURL, encodedPrompt)

	println("usedURL : ", usedURL)

	response, err := http.Get(usedURL)
	if err != nil {
		return nil, errors.New("cant send request to text-ai service")
	}

	body, err := io.ReadAll(response.Body)
   	if err != nil {
      return nil, errors.New("cant read response body")
   	}

	var responseData TextResponseData
	// แปลง body ให้เข้า format list ของ float64
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
