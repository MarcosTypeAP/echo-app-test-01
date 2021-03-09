package main

import (
	_ "database/sql"
	"fmt"
	"os"
)

//IsErr prints the error that was send to it
func IsErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func main() {

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	RunServer("", port)
}
