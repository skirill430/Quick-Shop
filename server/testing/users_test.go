package test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/skirill430/Quick-Shop/server/router"
	"github.com/skirill430/Quick-Shop/server/utils/db"
	"github.com/stretchr/testify/assert"
)

var Router = router.Router()

// needed to connect to DB before each test is run
func TestMain(m *testing.M) {
	db.ConnectDB("test")
	code := m.Run()

	// clear database after so future tests execute the same
	db.ClearDB()
	os.Exit(code)
}

func TestSignUp_OK(t *testing.T) {
	test_credentials := []byte(`{"username": "test-username", 
	"password": "test-password"}`)

	req, _ := http.NewRequest("POST", "/api/user/signup", bytes.NewBuffer(test_credentials))
	response := httptest.NewRecorder()
	Router.ServeHTTP(response, req)

	a := assert.New(t)
	a.Equal(http.MethodPost, req.Method, "HTTP request method error")
	a.Equal(http.StatusOK, response.Code, "HTTP request status code error")

	// delete newly created user in between tests, database is only cleared after all tests are executed
	db.DeleteUser("test-username")
}

func TestSignUp_TakenUsername(t *testing.T) {
	// Create first test user
	test_credentials := []byte(`{"username": "test-username", 
	"password": "test-password"}`)
	req, _ := http.NewRequest("POST", "/api/user/signup", bytes.NewBuffer(test_credentials))
	response := httptest.NewRecorder()
	Router.ServeHTTP(response, req)

	// Create duplicate second test user
	test_credentials = []byte(`{"username": "test-username", 
	"password": "test-password"}`)
	req, _ = http.NewRequest("POST", "/api/user/signup", bytes.NewBuffer(test_credentials))
	response = httptest.NewRecorder()
	Router.ServeHTTP(response, req)

	a := assert.New(t)
	a.Equal(http.MethodPost, req.Method, "HTTP request method error")
	a.Equal(http.StatusConflict, response.Code, "HTTP request status code error")

	db.DeleteUser("test-username")
}

func TestSignIn_OK(t *testing.T) {
	test_credentials := []byte(`{"username": "test-username", 
	"password": "test-password"}`)
	req, _ := http.NewRequest("POST", "/api/user/signup", bytes.NewBuffer(test_credentials))
	response := httptest.NewRecorder()
	Router.ServeHTTP(response, req)

	// attempt to sign in with newly created user with same credentials
	req, _ = http.NewRequest("POST", "/api/user/signin", bytes.NewBuffer(test_credentials))
	response = httptest.NewRecorder()
	Router.ServeHTTP(response, req)

	a := assert.New(t)
	a.Equal(http.MethodPost, req.Method, "HTTP request method error")
	a.Equal(http.StatusOK, response.Code, "HTTP request status code error")

	db.DeleteUser("test-username")
}

func TestSignIn_UsernameNotFound(t *testing.T) {
	// signing into user account 'test-username' that was never created
	test_credentials := []byte(`{"username": "test-username", 
	"password": "test-password"}`)
	req, _ := http.NewRequest("POST", "/api/user/signin", bytes.NewBuffer(test_credentials))
	response := httptest.NewRecorder()
	Router.ServeHTTP(response, req)

	a := assert.New(t)
	a.Equal(http.MethodPost, req.Method, "HTTP request method error")
	a.Equal(http.StatusNotFound, response.Code, "HTTP request status code error")
}

func TestSignIn_PasswordIncorrect(t *testing.T) {
	test_credentials := []byte(`{"username": "test-username", 
	"password": "test-password"}`)
	req, _ := http.NewRequest("POST", "/api/user/signup", bytes.NewBuffer(test_credentials))
	response := httptest.NewRecorder()
	Router.ServeHTTP(response, req)

	incorrect_credentials := []byte(`{"username": "test-username", 
	"password": "incorrect-password"}`)
	req, _ = http.NewRequest("POST", "/api/user/signin", bytes.NewBuffer(incorrect_credentials))
	response = httptest.NewRecorder()
	Router.ServeHTTP(response, req)

	a := assert.New(t)
	a.Equal(http.MethodPost, req.Method, "HTTP request method error")
	a.Equal(http.StatusUnauthorized, response.Code, "HTTP request status code error")

	db.DeleteUser("test-username")
}
