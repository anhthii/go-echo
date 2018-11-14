package models

import (
	"fmt"
	"os"

	"github.com/anhthii/go-echo/db"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type Token struct {
	Username string
	jwt.StandardClaims
}

// User store user's information
type User struct {
	Username     string `gorm:"PRIMARY_KEY"`
	Password     string
	Access_token string
	Playlists    []Playlist `gorm:"foreignkey:UserRefer"`
}

type errors map[string]string

func ValidateUsername(username string) (httpStatusCode int, errorResponse map[string]string) {
	user := &User{}
	if err := db.GetDB().Where("username = ?", username).First(&user).Error; gorm.IsRecordNotFoundError(err) {
		return 200, nil
	}

	return 400, errors{"username": fmt.Sprintf("%s username already in use by another user", user.Username)}
}

func CreateNewUser(username, password string) (httpStatusCode int, errorResponse map[string]string, dataResponse map[string]string) {
	user := &User{}
	user.Username = username
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)
	fmt.Printf("%v", user)

	db.GetDB().Create(user)

	// create new JWT token for the newly registered account
	tk := &Token{Username: user.Username}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("TOKEN_SECRET")))

	user.Access_token = tokenString
	db.GetDB().Save(&user)
	dataResponse = make(map[string]string)
	dataResponse["username"] = user.Username
	dataResponse["access_token"] = user.Access_token
	return 200, nil, dataResponse
}

func Login(username, password string) (httpStatusCode int, errorResponse map[string]string, dataResponse map[string]string) {
	user := &User{}
	if err := db.GetDB().Where("username = ?", username).First(&user).Error; gorm.IsRecordNotFoundError(err) {
		return 401, errors{"username": "Invalid username"}, nil
	}
	fmt.Println("token", user.Access_token)
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return 401, errors{"username": "Invalid password"}, nil
	}

	dataResponse = make(map[string]string)
	dataResponse["username"] = user.Username
	dataResponse["access_token"] = user.Access_token
	return 200, nil, dataResponse
}
