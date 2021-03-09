package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//IsErr prints the error that was send to it
func IsErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

//User user
type User struct {
	Name string `json:"name" xml:"name" form:"name" query:"name"`
	Age  int    `json:"age" xml:"age" form:"age" query:"age"`
}

func main() {
	e := echo.New()

	//e.Pre(middleware.HTTPSNonWWWRedirect())

	e.Use(middleware.CORS())

	e.Pre(middleware.AddTrailingSlash())

	e.Use(middleware.BasicAuth(auth))

	e.Use(middleware.Secure())

	config := middleware.RateLimiterConfig{
		Store: middleware.NewRateLimiterMemoryStore(5),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			id := ctx.RealIP()
			return id, nil
		},
	}

	e.Use(middleware.RateLimiterWithConfig(config))

	e.GET("/", home)

	e.Static("/googleclone", "./public")

	e.GET("/user/:id", getUser)
	e.POST("/user", postUser)
	e.POST("/user/json", postUserJSON)
	e.POST("/user/streamjson", postStreamJSON)
	// e.PUT("/user/:id", putUser)
	// e.DELETE("/user/:id", deleteUser)

	e.Logger.Fatal(e.Start(":1323"))
}

func auth(username, password string, c echo.Context) (bool, error) {
	if username == "user" && password == "passw" {
		return true, nil
	}
	return false, nil
}

func home(c echo.Context) error {
	return c.String(http.StatusOK, "Hello World!")
}

func getUser(c echo.Context) error {
	id := c.Param("id")
	return c.String(http.StatusOK, id)
}

func postUser(c echo.Context) error {
	name := c.FormValue("name")

	img, err := c.FormFile("image")
	IsErr(err)

	src, err := img.Open()
	IsErr(err)
	defer src.Close()

	dst, err := os.Create(img.Filename)
	IsErr(err)
	defer dst.Close()

	_, err = io.Copy(dst, src)
	IsErr(err)

	return c.String(http.StatusOK, name)
}

func postUserJSON(c echo.Context) error {
	var u User

	err := c.Bind(&u)
	IsErr(err)

	return c.JSON(http.StatusCreated, u)
}

func postStreamJSON(c echo.Context) error {
	u := &User{
		Name: "Jon",
		Age:  18,
	}
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response()).Encode(u)
}
