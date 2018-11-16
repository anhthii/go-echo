package user_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/anhthii/go-echo/db"
	"github.com/anhthii/go-echo/db/models"
	"github.com/anhthii/go-echo/router"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type UserTests struct{ Test *testing.T }

func TestRunner(t *testing.T) {
	db.Init("../../.env")
	db.Tables(&models.User{})
	defer db.Close()
	defer db.GetDB().DropTable("users")
	t.Run("Signup user", func(t *testing.T) {
		test := &UserTests{Test: t}
		test.ShouldRegisterUser()
		test.UserAlreadyExists()
		// test.TestCreateRegularUser()
		// test.TestCreateConfirmedUser()
		// test.TestCreateMasterUser()
		// test.TestCreateUserTwice()
	})
	t.Run("should login user", func(t *testing.T) {
		// test:= UserTests{Test: t}
		// test.TestLoginRegularUser()
		// test.TestLoginConfirmedUser()
		// test.TestLoginMasterUser()
	})

}

func (t *UserTests) ShouldRegisterUser() {

	postBody := map[string]string{
		"username": "testuser123",
		"password": "123456",
	}

	body, _ := json.Marshal(postBody)

	req, _ := http.NewRequest("POST", "/api/user/signup", bytes.NewReader(body))
	w := httptest.NewRecorder()
	r := router.SetupRouter("prod")

	r.ServeHTTP(w, req)
	assert.Equal(t.Test, w.Code, 200, "should get 200 response code ")
	var userResponse map[string]interface{}
	bytes, _ := ioutil.ReadAll(w.Body)
	json.Unmarshal(bytes, &userResponse)
	// check if response have username
	assert.Equal(t.Test, userResponse["username"], "testuser123", "username should be equal")
	// create jwt token from username
	accessToken := models.GenerateTokenForUser(postBody["username"])
	// check if response have access_token
	assert.Equal(t.Test, userResponse["access_token"], accessToken, "access_token should be the same")
}

func (t *UserTests) UserAlreadyExists() {

	postBody := map[string]string{
		"username": "testuser123",
		"password": "123456",
	}

	body, _ := json.Marshal(postBody)

	req, _ := http.NewRequest("POST", "/api/user/signup", bytes.NewReader(body))
	w := httptest.NewRecorder()
	r := router.SetupRouter("prod")

	r.ServeHTTP(w, req)
	assert.Equal(t.Test, w.Code, 400, "should get 400 response code")
	var userResponse map[string]interface{}
	bytes, _ := ioutil.ReadAll(w.Body)
	json.Unmarshal(bytes, &userResponse)

	responseErrors := userResponse["errors"].(map[string]interface{})
	// check if response have username
	assert.Equal(t.Test, responseErrors["username"].(string), postBody["username"]+" username already in use by another user", "username should be in use by an other user")
}
