# :computer: webserver in GO
A simple webserver written in GO using only the standard library package:

- Only accepts POST requests
(405 Method Not Allowed : A request method is not supported for the requested resource; for example, a GET request on a form that requires data to be presented via POST, or a PUT request on a read-only resource.)
- For any request, checks the body for a json encoded object and returns a descriptive HTTP error code if the body does not contain a valid json object.
	We expect this object to look like this: {"favoriteTree": "baobab"}
- Only accepts requests on the index URL: "/", and returns the proper HTTP error code if a different URL is requested.
- Runs locally on port 8000
- For a successful request, returns a properly encoded HTML document with the following content:
	If "favoriteTree" was specified:
		It's nice to know that your favorite tree is a <value of "favoriteTree" in the POST body>
	if not specified:
		Please tell me your favorite tree
		
# :whale: Run with Docker
```
docker-compose up
```

# :tv: View with Postman
https://www.getpostman.com/
 
#  :memo: Curl commands to test code:
(For windows need to escape " and use only double quotes.)

- Post request with correct body
```
curl -s -S -XPOST -d"{\"favoriteTree\":\"Beech\"}" http://localhost:8000
```

- Post request without body
```
curl -s -S -XPOST http://localhost:8000
```

- Post request with wrong body content
```
curl -s -S -XPOST -d"{\"name\":\"Beech\"}" http://localhost:8000
```

- Get request
```
curl http://localhost:8000
```
