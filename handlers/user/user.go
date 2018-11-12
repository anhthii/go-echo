package user

import (
	"fmt"
	"net/http"
	"strings"

	UserModel "github.com/anhthii/go-echo/db/models"
	"github.com/gin-gonic/gin"
	validator "gopkg.in/go-playground/validator.v8"
)

type UserBody struct {
	Username string `json:"username" binding:"required,alphanum,max=16,min=3"`
	Password string `json:"password" binding:"required,min=6"`
}

type ErrorResponse struct {
	Error  bool              `json:"error"` // define whether validate function return any errors or not
	Errors map[string]string `json:"errors"`
}

func validationErrorToText(e *validator.FieldError) string {
	switch e.Tag {
	case "required":
		return fmt.Sprintf("%s is is required", e.Field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", e.Field, e.Param)
	case "max":
		return fmt.Sprintf("%s must not be longer than %s characters", e.Field, e.Param)
	case "alphanum":
		return fmt.Sprintf("%s must be alphanumeric", e.Field)
	}
	return fmt.Sprintf("%s is not valid", e.Field)
}

func validate(errors validator.ValidationErrors) ErrorResponse {
	errorMap := make(map[string]string)
	haveError := false
	for _, err := range errors {
		errorMap[strings.ToLower(err.Field)] = validationErrorToText(err)
		haveError = true
	}

	return ErrorResponse{Error: haveError, Errors: errorMap}
}

// CreateNewUser check if data meets all the requirements, then create a new user
func CreateNewUser(c *gin.Context) {
	var json UserBody
	if err := c.ShouldBindJSON(&json); err != nil {
		validateResponse := validate(err.(validator.ValidationErrors))
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
		validateResponse := validate(err.(validator.ValidationErrors))
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
