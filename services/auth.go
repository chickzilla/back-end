package services

import (
	"database/sql"
	"net/http"
	"net/mail"

	"github.com/Her_feeling/back-end/database"
	utils "github.com/Her_feeling/back-end/utils/helper"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func SignUp(c *gin.Context) {
	var SignUpRequest AuthRequest
	var DB = database.DB

	// ถ้า body ที่เข้ามาไม่ตรงกับ user ให้ return 400
	if err := c.ShouldBindJSON(&SignUpRequest); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	// ตอนนี้ request มี email และ password แล้ว

	if _, err := mail.ParseAddress(SignUpRequest.Email); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	var foundEmail string

	err := DB.QueryRow(
		"SELECT email FROM user WHERE email = ?", SignUpRequest.Email).Scan(&foundEmail)

	if err != nil {
		if err == sql.ErrNoRows {
			hashed, err := bcrypt.GenerateFromPassword([]byte(SignUpRequest.Password), 8)

			if err != nil {
				c.JSON(http.StatusInternalServerError, err.Error())
				return
			}

			_, err = DB.Exec(
				"INSERT INTO user (email, password) VALUES (? , ?)",
				SignUpRequest.Email, hashed)

			if err != nil {
				c.JSON(http.StatusInternalServerError, err.Error())
				return
			} else {
				c.JSON(http.StatusCreated, gin.H{"response": "created successfully"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "email already exists"})
}

func SignIn(c *gin.Context) {
	var signupRequest AuthRequest
	DB := database.DB

	if err := c.ShouldBindJSON(&signupRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var hashedPassword string
	err := DB.QueryRow("SELECT password FROM user WHERE email = ?", signupRequest.Email).Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "email not found"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(signupRequest.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong password"})
		return
	}

	jwtToken, err := utils.GenerateKey(signupRequest.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": jwtToken})

}
