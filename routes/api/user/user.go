package user

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	validator "gopkg.in/go-playground/validator.v8"
)

type User struct {
	Username string `json:"username" binding:"required,alphanum,max=6,min=3"`
	Password string `json:"password" binding:"required,min=6"`
}

type ValidateResponse struct {
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

func validate(errors validator.ValidationErrors) ValidateResponse {
	errorMap := make(map[string]string)
	haveError := false
	for _, err := range errors {
		errorMap[strings.ToLower(err.Field)] = validationErrorToText(err)
		haveError = true
	}
	return ValidateResponse{Error: haveError, Errors: errorMap}
}

// CreateNewUser check if data meets all the requirements, then create a new user
func CreateNewUser(c *gin.Context) {
	var json User
	if err := c.ShouldBindJSON(&json); err != nil {
		validateResponse := validate(err.(validator.ValidationErrors))
		c.JSON(http.StatusBadRequest, validateResponse)
		return
	}
	if json.Username != "manu" || json.Password != "123456" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
}
