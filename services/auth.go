package services

import (
	"database/sql"
	"fmt"
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

	if err := c.ShouldBindJSON(&SignUpRequest); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

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
	fmt.Println("new version none mode 2")

	var signupRequest AuthRequest
	DB := database.DB

	if err := c.ShouldBindJSON(&signupRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var hashedPassword *string
	var onlySSO bool

	err := DB.QueryRow("SELECT password, only_SSO FROM user WHERE email = ?", signupRequest.Email).Scan(&hashedPassword, &onlySSO)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "email not found"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if onlySSO {
		c.JSON(http.StatusBadRequest, gin.H{"error": "this email can access only by SSO"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(*hashedPassword), []byte(signupRequest.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong password"})
		return
	}

	jwtToken, err := utils.GenerateKey(signupRequest.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	/*
		domain, _ := utils.GetEnvNoCon("DOMAIN_URL")
		c.SetSameSite(http.SameSiteNoneMode)
		c.SetCookie("auth_token", jwtToken, 3600*24*30, "/", domain, true, false)*/

	c.JSON(http.StatusOK, gin.H{"response": jwtToken})

}

type AuthSSORequest struct {
	Email string `json:"email" binding:"required"`
}

func SignInWithSSO(c *gin.Context) {

	fmt.Println("new version none mode 2")

	var SSORequest AuthSSORequest
	DB := database.DB

	if err := c.ShouldBindBodyWithJSON(&SSORequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var foundEmail string

	err := DB.QueryRow("SELECT email FROM user WHERE email = ?", SSORequest.Email).Scan(&foundEmail)
	if err != nil {
		if err == sql.ErrNoRows {

			_, err = DB.Exec("INSERT INTO user (email, only_SSO) VALUES (?, ?)", SSORequest.Email, true)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	jwtToken, err := utils.GenerateKey(SSORequest.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	/*
		domain, _ := utils.GetEnvNoCon("DOMAIN_URL")
		c.SetCookie("auth_token", jwtToken, 3600*24*30, "/", domain, true, false)
	*/
	c.JSON(http.StatusOK, gin.H{"response": jwtToken})
}

func SignOut(c *gin.Context) {
	fmt.Println("new version none mode")

	/*
		domain, _ := utils.GetEnvNoCon("DOMAIN_URL")
		c.SetCookie("auth_token", "", -1, "/", domain, false, true)
	*/

	c.JSON(http.StatusOK, gin.H{"message": "Signed out"})

}
