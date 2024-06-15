package utils

import (
	"errors"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

func GetEnv(wg *sync.WaitGroup, envChan chan<- string, errChan chan<- error) {
	defer wg.Done() 

	err := godotenv.Load()
	if err != nil {
		errChan <- errors.New("error while loading .env file")
		return
	}

	textAIURL := os.Getenv("AI_TEXT_URL")
		if textAIURL == ""{
			errChan <- errors.New("cannot read AI_TEXT_URL in .env file")
			return
		}


	envChan <- textAIURL
}