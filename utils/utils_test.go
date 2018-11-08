package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const URL = "https://jsonplaceholder.typicode.com/todos/1"

func TestGetMapDataFromHTTPGet(t *testing.T) {
	data, err := GetMapDataFromHTTPGet(URL)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, int(data["userId"].(float64)), 1, "userId should be equal")
	assert.Equal(t, int(data["id"].(float64)), 1, "id should be equal")
	assert.Equal(t, data["title"].(string), "delectus aut autem", "title should be equal")
	assert.Equal(t, data["completed"].(bool), false, "completed should be as expected")
}
