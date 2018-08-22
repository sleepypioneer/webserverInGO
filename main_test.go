package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Create a ResponseRecorder (which satisfies http.ResponseWriter) to record the test responses.
func createResponseRecorder() (*httptest.ResponseRecorder, http.Handler) {
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(fromIndex(postRequest(requestHandler)))
	return rr, handler
}

// Test with POST method but not from root URL
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
	d := make(map[string]interface{})
	d["name"] = "Jessica"
	b, err := json.Marshal(d)
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
	d := make(map[string]interface{})
	d["name"] = "Jessica"
	b, err := json.Marshal(d)
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
	// Check the response body contains the given tree
	if name := "Jessica"; !strings.Contains(rr.Body.String(), name) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), name)
	}
}

// Check response body with POST method and incorrect body
func Test_responseHandler_4(t *testing.T) {
	d := make(map[string]interface{})
	d["userName"] = "Jessica"
	b, err := json.Marshal(d)
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
			rr.Body.String(), expected)
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
