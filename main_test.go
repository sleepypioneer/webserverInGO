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

/* // Create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
rr := httptest.NewRecorder()
handler := http.HandlerFunc(fromIndex(postRequest(requestHandler)))
handler.ServeHTTP(rr, req) */

// Test with POST method but not root URL
func Test_responseHandler_1(t *testing.T) {
	req, err := http.NewRequest("POST", "/smth", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(fromIndex(postRequest(requestHandler)))
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

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(fromIndex(postRequest(requestHandler)))
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

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(fromIndex(postRequest(requestHandler)))
	handler.ServeHTTP(rr, req)

	// Check the response body is HTML
	expected := "<!DOCTYPE html>"
	if !strings.Contains(rr.Body.String(), "<!DOCTYPE html>") {
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

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(fromIndex(postRequest(requestHandler)))
	handler.ServeHTTP(rr, req)

	// Check the returned document contains the expected string
	expected := "Please tell me your favorite tree?"
	if !strings.Contains(rr.Body.String(), "Please tell me your favorite tree?") {
		t.Errorf("handler returned unexpected body: got %v want %v",
			len(rr.Body.String()), expected)
	}
}

// Test with POST method and from root with no body
func Test_responseHandler_5(t *testing.T) {
	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(fromIndex(postRequest(requestHandler)))
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusPreconditionFailed {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusPreconditionFailed)
	}
}
