package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.GET("/", f)

	e.Logger.Fatal(e.Start(":1323"))
}

func f(c echo.Context) error {
	return c.String(http.StatusOK, "Hello World!")
}
