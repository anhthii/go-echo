package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	validator "gopkg.in/go-playground/validator.v8"
)

func httpGet(URL string) ([]byte, error) {
	response, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// GetMapDataFromHTTPGet return data with a map structure read from http get response
func GetMapDataFromHTTPGet(URL string) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	bytes, err := httpGet(URL)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(bytes, &data)
	return data, nil
}

// GetStringDataFromHTTPGet return data as a string read from http get response
func GetStringDataFromHTTPGet(URL string) (string, error) {
	bytes, err := httpGet(URL)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// InternalErrorJSON return internal errors as JSON response
func InternalErrorJSON(context *gin.Context, err error) {
	_error := gin.H{"errors": err}
	context.JSON(500, _error)
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

func Validate(errors validator.ValidationErrors) ErrorResponse {
	errorMap := make(map[string]string)
	haveError := false
	for _, err := range errors {
		errorMap[strings.ToLower(err.Field)] = validationErrorToText(err)
		haveError = true
	}

	return ErrorResponse{Error: haveError, Errors: errorMap}
}
