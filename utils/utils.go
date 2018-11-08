package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
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
