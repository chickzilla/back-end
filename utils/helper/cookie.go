package utils

import (
	"github.com/gin-gonic/gin"
)

func GetCookie(c *gin.Context, tokenName string) (string, error) {
	cookie, err := c.Cookie(tokenName)
	if err != nil {
		return "", err
	}

	return cookie, nil
}
