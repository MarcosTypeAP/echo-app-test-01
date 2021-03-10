package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//RunServer runs the server
func RunServer(address string, port string) {
	e := echo.New()

	e.HideBanner = true

	// e.Pre(middleware.AddTrailingSlash())
	e.Use(middleware.CORS())

	server := &http.Server{
		Addr:           ":" + port,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	e.POST("/test", func(c echo.Context) error {
		username := c.FormValue("username")
		password := c.FormValue("password")
		fmt.Println(username, password)
		return c.NoContent(200)
	})

	e.GET("/", homeHandler)

	e.Static("/g", "./pages/google-clone")

	e.POST("/signup", signUpHandler)
	e.Static("/login", "./pages/log-in")
	e.GET("/unauthorized", unauthorizedPageHandler)

	n := e.Group("/notesapp", logInMiddleware)
	n.GET("/notes", userNotesPageHandler)
	n.POST("/notes", userNotesPageHandler)
	e.File("/notesapp/main.css", "./pages/notes/main.css")

	e.POST("/loginjwt", logInWithJWTHandler, logInWithJWTMiddleware)

	e.Logger.Fatal(e.StartServer(server))
}

func logInWithJWTHandler(c echo.Context) error {
	return c.String(200, "sarpado")
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

func unauthorizedPageHandler(c echo.Context) error {
	return c.String(http.StatusUnauthorized, "Hasta aca llegaste pa")
}

func userNotesPageHandler(c echo.Context) error {
	return c.File("./pages/notes/index.html")
}

func logInMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		username := c.FormValue("username")
		password := c.FormValue("password")
		// fmt.Println(username)

		if username == "" {
			return c.Redirect(http.StatusMovedPermanently, "/login/")
		}

		hashedPassword, err := GetUserPasswordDB(username)
		IsErr(err)
		if err != nil {
			return c.Redirect(http.StatusMovedPermanently, "/unauthorized")
		}

		match, err := ComparePassword(password, hashedPassword)
		IsErr(err)

		if !match {
			return c.Redirect(http.StatusMovedPermanently, "/unauthorized")
		}

		return next(c)
	}
}

func logInWithJWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		username := c.FormValue("username")
		password := c.FormValue("password")
		// fmt.Println(username)

		if username == "" {
			return c.Redirect(http.StatusMovedPermanently, "/login/")
		}

		hashedPassword, err := GetUserPasswordDB(username)
		IsErr(err)
		if err != nil {
			return c.Redirect(http.StatusMovedPermanently, "/unauthorized")
		}

		match, err := ComparePassword(password, hashedPassword)
		IsErr(err)

		if !match {
			return c.Redirect(http.StatusMovedPermanently, "/unauthorized")
		}

		token, err := GenerateJWT(username, "admin")
		if err != nil {
			return c.String(500, err.Error())
		}

		c.JSON(200, map[string]string{"token": token})
		fmt.Println(token)

		return next(c)
	}
}
