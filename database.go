package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

//User is an user
type User struct {
	UserID   int    `json:"user_id,omitempty"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
}

var connectionString string = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
	os.Getenv("user"),
	os.Getenv("pass"),
	os.Getenv("host"),
	os.Getenv("db_port"),
	os.Getenv("db_name"),
)

// OpenConnectionDB opens connection to the database
func OpenConnectionDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// CheckUsernameExists checks if the username exists
func CheckUsernameExistsDB(username string) (bool, error) {
	db, err := OpenConnectionDB()
	defer db.Close()
	if err != nil {
		return false, err
	}

	result, err := db.Query("SELECT true FROM `user` WHERE `username` = '" + username + "'")
	defer result.Close()
	if err != nil {
		return false, err
	}

	var exist bool

	result.Next()
	_ = result.Scan(&exist)

	return exist, err
}

// AddUser adds a new user to the database
func AddUserDB(username string, hashedPassword []byte) error {
	db, err := OpenConnectionDB()
	defer db.Close()
	if err != nil {
		return err
	}

	insert, err := db.Query(fmt.Sprintf("INSERT INTO `user` (`username`, `password`) VALUES('%s', '%s')",
		username,
		string(hashedPassword),
	))
	defer insert.Close()
	if err != nil {
		return err
	}

	return nil
}

// GetUserPasswordDB gets the hashed password of the searched user
func GetUserPasswordDB(username string) (string, error) {
	db, err := OpenConnectionDB()
	defer db.Close()
	if err != nil {
		return "", err
	}

	result, err := db.Query("SELECT `password` FROM `user` WHERE `username` = '" + username + "'")
	defer result.Close()
	if err != nil {
		return "", err
	}

	var hashedPassword string

	result.Next()
	err = result.Scan(&hashedPassword)

	return hashedPassword, err
}
