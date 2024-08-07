package database

import (
	"database/sql"
	"fmt"

	utils "github.com/Her_feeling/back-end/utils/helper"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	dbURL, err := utils.GetEnvNoCon("MYSQL_URL")
	if err != nil {
		panic("cant get mysql_url in .env file")
	}

	db, err := sql.Open("mysql", dbURL)
	if err != nil {
		panic("Database could not connect, error : " + err.Error())
	}

	fmt.Println("Connect to database successfully!")

	DB = db
	createTables()
}

func createTables() error {
	createdUserTable := `
	CREATE TABLE IF NOT EXISTS user (
		id INT AUTO_INCREMENT PRIMARY KEY,
		email VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255),
		only_SSO BOOLEAN NOT NULL DEFAULT FALSE,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)
	`

	_, err := DB.Exec(createdUserTable)

	if err != nil {
		panic("Could not create user table, error : " + err.Error())
	}

	createPromptResultLogsTable := `
	CREATE TABLE IF NOT EXISTS user_history (
		id INT AUTO_INCREMENT PRIMARY KEY,
		user_id INT NOT NULL,
		prompt TEXT NOT NULL,
		love_prob FLOAT(2),
		sadness_prob FLOAT(2),
		joy_prob FLOAT(2),
		angry_prob FLOAT(2),
		fear_prob FLOAT(2),
		surprise_prob FLOAT(2),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(user_id) REFERENCES user(id)
	)
	`

	_, err = DB.Exec(createPromptResultLogsTable)

	if err != nil {
		panic("Could not create user_history table, error : " + err.Error())
	}

	fmt.Println("Init table successfully!")

	return err

}
