package user

import (
	"net/http"

	UserModel "github.com/anhthii/go-echo/db/models"
	. "github.com/anhthii/go-echo/utils"
	"github.com/gin-gonic/gin"
	validator "gopkg.in/go-playground/validator.v8"
)

type UserBody struct {
	Username string `json:"username" binding:"required,alphanum,max=16,min=3"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginBody struct {
	UserBody
}

// CreateNewUser check if data meets all the requirements, then create a new user
func CreateNewUser(c *gin.Context) {
	var json UserBody
	if err := c.ShouldBindJSON(&json); err != nil {
		validateResponse := Validate(err.(validator.ValidationErrors))
		c.JSON(http.StatusBadRequest, validateResponse)
		return
	}

	httpStatusCode, errorResponse := UserModel.ValidateUsername(json.Username)
	if errorResponse != nil {
		c.JSON(httpStatusCode, ErrorResponse{Error: true, Errors: errorResponse})
		return
	}

	httpStatusCode, errorResponse, dataResponse := UserModel.CreateNewUser(json.Username, json.Password)
	if errorResponse != nil {
		c.JSON(httpStatusCode, ErrorResponse{Error: true, Errors: errorResponse})
		return
	}

	c.JSON(http.StatusOK, dataResponse)
}

func Login(c *gin.Context) {
	var json UserBody
	if err := c.ShouldBindJSON(&json); err != nil {
		validateResponse := Validate(err.(validator.ValidationErrors))
		c.JSON(http.StatusBadRequest, validateResponse)
		return
	}
	httpStatusCode, errorResponse, dataResponse := UserModel.Login(json.Username, json.Password)
	if errorResponse != nil {
		c.JSON(httpStatusCode, ErrorResponse{Error: true, Errors: errorResponse})
		return
	}

	c.JSON(http.StatusOK, dataResponse)
}
