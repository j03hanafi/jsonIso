package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Init(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()

	server().ServeHTTP(response, request)
	expectedResponse := `{"message":"Hello World!","status":200}`
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Log(err)
	}
	assert.Equal(t, 200, response.Code, "Invalid response code")
	assert.Equal(t, expectedResponse, string(bytes.TrimSpace(responseBody)))
	t.Log(response.Body.String())
}

func Test_GetPayments(t *testing.T) {
	request, _ := http.NewRequest("GET", "/payment", nil)
	response := httptest.NewRecorder()

	server().ServeHTTP(response, request)
}
