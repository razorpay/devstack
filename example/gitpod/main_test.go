package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
//	"strings"
	"github.com/stretchr/testify/assert"
)

func TestHelloword(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(helloworld)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method 
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	expected_status := http.StatusOK
	assert.Equal(t,rr.Code,expected_status,"handler returned wrong status code: got %v want %v",rr.Code, expected_status)
	
	expected_text := "Helloworld"
	assert.HTTPBodyContains(t, helloworld, "GET", "/", nil, expected_text ,"handler returned wrong status code: got %v want %v",rr.Code, expected_text)
}