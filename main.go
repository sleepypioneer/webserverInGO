/*
A simple web server, using only the standard library package of go:

* Only accepts POST requests
(405 Method Not Allowed : A request method is not supported for the requested resource; for example, a GET request on a form that requires data to be presented via POST, or a PUT request on a read-only resource.)
* For any request, checks the body for a json encoded object and returns a descriptive HTTP error code if the body does not contain a valid json object.
	We expect this object to look like this: {"favoriteTree": "baobab"}
* Only accepts requests on the index URL: "/", and returns the proper HTTP error code if a different URL is requested.
* Runs locally on port 8000
* For a successful request, returns a properly encoded HTML document with the following content:
	If "favoriteTree" was specified:
		It's nice to know that your favorite tree is a <value of "favoriteTree" in the POST body>
	if not specified:
		Please tell me your favorite tree
*/

package main

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
)

// Log requests
func logging(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("---New Request----")
		log.Println("Path: ", r.URL.Path)
		log.Println("Method: ", r.Method)
		log.Println("Body: ", r.Body)
		log.Println("------------------")
		f.ServeHTTP(w, r)
	}
}

// Only allow requests from root
func fromIndex(f http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		f.ServeHTTP(w, r)
	})
}

// Only allow POST method requests
func postRequest(f http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Please only use Post requests", http.StatusMethodNotAllowed)
			return
		}
		f.ServeHTTP(w, r)
	})
}

// Handle request and decode JSON data if found
func requestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Body != nil {
		d := make(map[string]interface{})
		err := json.NewDecoder(r.Body).Decode(&d)
		if err != nil {
			if err == io.EOF {
				http.Error(w, "Body of request cannot be empty, expecting Json data.", http.StatusPreconditionFailed)
				return
			}
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		responseHandler(w, d)
		return
	}
	http.Error(w, "Body of request cannot be nil, expecting Json data.", http.StatusTeapot)
}

// Handle response depedning on if favoriteTree has been specified
func responseHandler(w http.ResponseWriter, d map[string]interface{}) {
	t, err := template.ParseFiles("./templates/favoriteTree.html")
	nt, err := template.ParseFiles("./templates/noTree.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	v, ok := d["favoriteTree"]
	if ok {
		err := t.Execute(w, v)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		err := nt.Execute(w, v)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
}

func main() {
	http.HandleFunc("/", logging(fromIndex(postRequest(requestHandler)))) // set router
	err := http.ListenAndServe(":8000", nil)                              // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
