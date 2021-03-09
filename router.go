package main

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

//RunServer runs the server
func RunServer(address string, port string) {
	e := echo.New()

	e.HideBanner = true

	server := &http.Server{
		Addr:           ":" + port,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	e.GET("/", homeHandler)

	e.POST("/signup", signUpHandler)
	e.GET("/login", logInHandler)
	e.GET("/notes", userNotesPageHandler, logInMiddleware)

	e.Logger.Fatal(e.StartServer(server))
}

func homeHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Home")
}

func signUpHandler(c echo.Context) error {

	username := c.FormValue("username")
	password := c.FormValue("password")

	err := VerifyIfUsernameIsValid(username)
	IsErr(err)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotAcceptable, err.Error())
	}
	err = VerifyIfPasswordIsValid(password)
	IsErr(err)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotAcceptable, err.Error())
	}

	hashedPassword, err := HashPassword(password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = AddUserDB(username, hashedPassword)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusCreated, "user created")
}

func logInHandler(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	hashedPassword, err := GetUserPasswordDB(username)
	IsErr(err)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
	}

	match, err := ComparePassword(password, hashedPassword)
	IsErr(err)

	if !match {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
	}

	// userNotesPageHandler(c)
	// return c.NoContent(http.StatusMovedPermanently)
	return c.Redirect(http.StatusMovedPermanently, "/")
}

func unauthorizedPageHandler(c echo.Context) error {
	return c.String(http.StatusUnauthorized, "Hasta aca llegaste pa")
}

func userNotesPageHandler(c echo.Context) error {
	return c.String(200, "en un futuro habrá algo aquí xd")
}

func logInMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		username := c.FormValue("username")
		password := c.FormValue("password")

		hashedPassword, err := GetUserPasswordDB(username)
		IsErr(err)
		if err != nil {
			return unauthorizedPageHandler(c)
			// return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
		}

		match, err := ComparePassword(password, hashedPassword)
		IsErr(err)

		if !match {
			return unauthorizedPageHandler(c)
			// return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
		}

		return next(c)
	}

}
