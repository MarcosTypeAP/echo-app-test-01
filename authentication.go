package main

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"io/ioutil"
	"time"
	"unicode"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type jwtCustomClaims struct {
	User string `json:"user"`
	Role string `json:"role"`
	jwt.StandardClaims
}

var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
)

func init() {
	privateKeyBytes, err := ioutil.ReadFile("./keys/private.rsa.key")
	IsErr(err)

	publicKeyBytes, err := ioutil.ReadFile("./keys/public.rsa.pub")
	IsErr(err)

	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	IsErr(err)

	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	IsErr(err)
}

// GenerateJWT generates a JWT
func GenerateJWT(username, role string) (string, error) {

	claims := jwtCustomClaims{
		username,
		role,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	t, err := token.SignedString(privateKey)
	IsErr(err)
	if err != nil {
		return "", err
	}

	return t, nil
}

func g(c echo.Context) {
	token, err := jwt.ParseFromRequestWithClaims(c.Request, request.OAuth2Extractor, &jwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})

	fmt.Println(token, err.Error())
}

// HashPassword hashs a password with cost 10
func HashPassword(password string) ([]byte, error) {
	passInBytes := []byte(password)

	hashedPassword, err := bcrypt.GenerateFromPassword(passInBytes, 10) //DefaultCost is 10
	IsErr(err)

	return hashedPassword, err
}

// ComparePassword compares a password with its password hashed
func ComparePassword(password, hashedPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	IsErr(err)

	if err == nil {
		return true, nil
	}

	return false, err
}

// VerifyIfPasswordIsValid verifies if the password follow the rules to be valid
func VerifyIfPasswordIsValid(password string) error {
	if len(password) < 8 {
		return errors.New("password is less than 8 characters")
	}

	validate := validator.New()

	if err := validate.Var(password, "printascii"); err != nil {
		return err
	}

	var (
		upp, low, num, sym bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char) && upp == false:
			upp = true
		case unicode.IsLower(char) && low == false:
			low = true
		case unicode.IsNumber(char) && num == false:
			num = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char) && sym == false:
			sym = true
		}
	}

	switch {
	case !upp:
		return errors.New("password does not contain uppercase")
	case !low:
		return errors.New("password does not contain lowercase")
	case !num:
		return errors.New("password does not contain numbers")
	case !sym:
		return errors.New("password does not contain special characters")
	default:
		return nil
	}
}

// VerifyIfUsernameIsValid verifies if the password follow the rules to be valid
func VerifyIfUsernameIsValid(username string) error {
	if len(username) < 4 {
		return errors.New("username is less than 4 characters")
	}

	validate := validator.New()

	err := validate.Var(username, "alphanum")
	if err != nil {
		return err
	}

	exist, err := CheckUsernameExistsDB(username)
	if err != nil {
		return err
	}
	if exist {
		return errors.New("username already exists")
	}

	return nil
}
