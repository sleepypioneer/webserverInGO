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
	"testing"
)

type testBody struct {
	FavoriteTree string
}

type testFalseBody struct {
	Name string
}

func Test_responseHandler_1(t *testing.T) {
	// Test with POST method but not root URL
	req, err := http.NewRequest("POST", "/smth", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(responseHandler)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func Test_responseHandler_2(t *testing.T) {
	// Test with POST method and correct body
	tb := testBody{FavoriteTree: "Beech"}
	b, err := json.Marshal(tb)
	r := bytes.NewReader(b)
	req, err := http.NewRequest("POST", "/", r)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(responseHandler)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func Test_responseHandler_3(t *testing.T) {
	// Check response body with POST method and correct body
	tb := testBody{FavoriteTree: "Oak"}
	b, err := json.Marshal(tb)
	r := bytes.NewReader(b)
	req, err := http.NewRequest("POST", "/", r)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(responseHandler)

	handler.ServeHTTP(rr, req)

	// Check the response body is what we expect by length.
	expected := 1589
	if len(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			len(rr.Body.String()), expected)
	}
}

func Test_responseHandler_4(t *testing.T) {
	// Check response body with POST method and incorrect body
	tb := testFalseBody{Name: "Jessica"}
	b, err := json.Marshal(tb)
	r := bytes.NewReader(b)
	req, err := http.NewRequest("POST", "/", r)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(responseHandler)

	handler.ServeHTTP(rr, req)

	// Check the response body is what we expect by length.
	expected := 1714
	if len(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			len(rr.Body.String()), expected)
	}
}

// Not currently working, goes into Panic dur to JSON stream reference
func Test_responseHandler_5(t *testing.T) {
	// Test with POST method and from root with no body
	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(responseHandler)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusPreconditionFailed {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusPreconditionFailed)
	}
}
