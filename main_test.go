/*
TEST CASES:
* Non-POST requests return 405 (Method Not Allowed)
* No JSON found in body returns a descriptive HTTP error code
* Check objest return matches this {"favoriteTree": "baobab"}
* Non-index URL: "/", returnsproper HTTP error code
* Check that for a successful request, returns a properly encoded HTML document with the following content:
	If "favoriteTree" was specified:
		It's nice to know that your favorite tree is a <value of "favoriteTree" in the POST body>
	if not specified:
		Please tell me your favorite tree
*/

package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type (
	testBody struct {
		FavoriteTree string
	}

	testFalseBody struct {
		Name string
	}
)

// Create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
func createResponseRecorder() (*httptest.ResponseRecorder, http.Handler) {
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(fromIndex(postRequest(requestHandler)))
	return rr, handler
}

// Test with POST method but not root URL
func Test_responseHandler_1(t *testing.T) {
	req, err := http.NewRequest("POST", "/smth", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr, handler := createResponseRecorder()
	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

// Test with POST method and correct body
func Test_responseHandler_2(t *testing.T) {
	tb := testBody{FavoriteTree: "Beech"}
	b, err := json.Marshal(tb)
	r := bytes.NewReader(b)
	req, err := http.NewRequest("POST", "/", r)
	if err != nil {
		t.Fatal(err)
	}
	rr, handler := createResponseRecorder()
	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

// Check response body with POST method and correct body
func Test_responseHandler_3(t *testing.T) {
	tb := testBody{FavoriteTree: "Oak"}
	b, err := json.Marshal(tb)
	r := bytes.NewReader(b)
	req, err := http.NewRequest("POST", "/", r)
	if err != nil {
		t.Fatal(err)
	}
	rr, handler := createResponseRecorder()
	handler.ServeHTTP(rr, req)
	// Check the response body is HTML
	if expected := "<!DOCTYPE html>"; !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

// Check response body with POST method and incorrect body
func Test_responseHandler_4(t *testing.T) {
	tb := testFalseBody{Name: "Jessica"}
	b, err := json.Marshal(tb)
	r := bytes.NewReader(b)
	req, err := http.NewRequest("POST", "/", r)
	if err != nil {
		t.Fatal(err)
	}
	rr, handler := createResponseRecorder()
	handler.ServeHTTP(rr, req)
	// Check the returned document contains the expected string
	if expected := "Please tell me your favorite tree?"; !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			len(rr.Body.String()), expected)
	}
}

// Check response body with POST method and empty body
func Test_responseHandler_5(t *testing.T) {
	body := bytes.NewReader([]byte{})
	req, err := http.NewRequest("POST", "/", body)
	if err != nil {
		t.Fatal(err)
	}
	rr, handler := createResponseRecorder()
	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusPreconditionFailed {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusPreconditionFailed)
	}
}

// Test with POST method and from root with nil body
func Test_responseHandler_6(t *testing.T) {
	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr, handler := createResponseRecorder()
	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusTeapot {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusTeapot)
	}
}
