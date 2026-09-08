# Backend Web Development Using Go
__Contents__:
- [__HTTP Clients__](#http-clients)
    - [JSON](#json)
    - [DNS](#dns)
- [__HTTP Servers__]
- [__Web Security__]

## HTTP Clients
Computers need to speak the same language (protocol) to communicate over the internet. The most used protocol is __HTTP__ (Hypertext Transfer Protocol). It is an __application-layer__ network protocol that defines a standardized, plain-text format for client and server computers to send requests and receive responses over the internet.

Benefits of Using HTTP:
- __Language Agnostic__: Request and responses can be sent from any program written in different languages.
- __Statelessness__: Each HTTP request is completely independent and contains all the information necessary to complete it. The server doesn't need to hold open a permanent, memory-heavy connection for every client on earth.
- __Structured and Human Readable__: HTTP messages are in plain-text with distinct components:
    - __Method__: Tells the server what action to take.
    - __URL/Path__: Specifies which resource is being requested.
    - __Headers__: key-value pairs providing metadata.
    - __Body__: The actual payload, normally in JSON.

A __Uniform Resource Locator (URL)__ is the address of another computer or server on the internet, with part of it specifying where to reach the server and the other part stating the information needed, such as `https//:www.hello.com/greet`. The `http` indicates that the __http__ protocol is being used.

In the request and response model the client is the computer that asks for information, while the server is the computer that responds by sending information. The `net/http` standard package allows requests to be made with methods like `http.Client` and `http.Get` which uses `http.defaultClient` under the hood.
| Function / Method | Example Code | Why It's Important / Primary Use Case |
| :--- | :--- | :--- |
| **`http.Get`** | `resp, err := http.Get("https://api.example.com/users")` | **Quick Reads:** Fetches data from a URL using an HTTP `GET` request. Best for simple queries where default client settings are sufficient. |
| **`http.Post`** | `resp, err := http.Post(url, "application/json", bodyReader)` | **Data Creation:** Sends data (like a JSON payload) to a server using a `POST` request to create new resources. |
| **`http.PostForm`** | `resp, err := http.PostForm(url, url.Values{"key": {"val"}})` | **Form Submissions:** Sends HTML form data (`application/x-www-form-urlencoded`) natively without manually formatting headers. |
| **`http.NewRequest`** | `req, err := http.NewRequest("PUT", url, bodyReader)` | **Custom Requests:** Creates a customizable request object (`*http.Request`). Required for HTTP verbs like `PUT`, `DELETE`, or `PATCH`, or when adding custom Headers (e.g., Auth tokens). |
| **`client.Do`** | `resp, err := client.Do(req)` | **Execution Engine:** Executes an `*http.Request` created by `NewRequest`. Lets you reuse custom `http.Client` configurations (timeouts, redirects). |
| **`http.Client{...}`** | `client := &http.Client{Timeout: 10 * time.Second}` | **Production Safety:** Configures client behavior. Crucial for setting global timeouts so your Go app doesn't hang indefinitely on slow responses. |
| **`resp.Body.Close`** | `defer resp.Body.Close()` | **Resource Cleanup:** Closes the response body network stream. **Mandatory** to prevent memory/connection leaks in Go network applications. |

__Important Mechanics__:
- Creating Custom Requests: Use `http.NewRequest` with `client.Do` because `http.Get` and `http.Post` use the unconfigured `httpDefaultClient` under the hood which has no timeout.
- Always Defer `resp.Body.Close`: Whenever `http.Get()`, `http.Post()`, or `client.Do()` succeeds (when `err == nil`), the Go runtime opens an active network connection to stream the response data. If you forget to call `resp.Body.Close()`, that connection is never returned to the system pool. If your server processes thousands of requests, it will quickly run out of socket file descriptors and crash!

```go
package main

import (
	"net/http"
	"time"
)

func main() {
	// 1. Create a safe client with a 5-second timeout
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// 2. Build a custom request (needed for PUT, DELETE, or setting Headers)
	req, err := http.NewRequest("DELETE", "https://api.example.com/users/123", nil)
	if err != nil {
		return
	}

	// 3. Attach custom metadata / auth headers
	req.Header.Set("Authorization", "Bearer MY_SECRET_TOKEN")

	// 4. Execute the request
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close() // ALWAYS close the body!
}
```
### JSON
__JavaScript Object Notation__ is the standard for transferring data across the web with the HTTP protocol. It is a __stringified__ version Javascript's object syntax and supports types like:
- __Primitive__: String, Numbers, Booleans, and Null
- __Collections__: Arrays, Object literals like dictionaries

It is commonaly used in:
- __HTTP__: in the body of the request and response objects.
- __Text files__: as configuration file (`.json`).
- __DBs__: In NoSQL databases like MongoDB, ElasticSearch, and Firestore

With JSON the keys are always strings, whereas the values can be any type it supports.
```json
{
    "movies": [
        {
            "id": 1,
            "title": "Iron Man",
            "director": "Jon Favreau",
            "favorite": true
        },
        {
            "id": 2,
            "title": "The Avengers",
            "director": "Joss Whedon",
            "favorite": false
        }
    ]
}
```
Using backticks in Go makes in easier to write JSON data through string literals.

JSON data is sent in the body of an HTTP response in the form of __bytes__ which need to be __decoded__. In Go, a __struct__ is used to decode the the data which requires knowing the JSON fields and their types. The `encoding/json` package provides a _compile time type-safety_ way to do so. To tell the translator how Go struct fields map to JSON keys, Go uses Struct Tags—metadata annotations written inside backticks `json:"..."`. The struct fields must be __capitalised__ and exported for the `encoding/json` package to access them. Use struct tags to map uppercase Go field names to lowercase JSON keys.
```go
type User struct {
	ID        int    `json:"id"`                 // Maps 'ID' to "id" in JSON
	Name      string `json:"full_name"`          // Maps 'Name' to "full_name"
	Email     string `json:"email,omitempty"`    // Omits field if empty/zero-value
	Password  string `json:"-"`                  // Ignores this field completely
}
```
When the response is received the values can be decoded into a _slice_ using the `&` _address of operator_.
```go
// The response is successful: http.Response
var users []User
decoder := json.NewDecoder(res.Body)
if err := decoder.Decode(&users); err != nil {
    fmt.Println("error decoding response")
    return
}
// If no error occurs the slice can be used
for _, user := range users {
    fmt.Printf("User - id: %v, name: %v, email: %v", user.ID, user.Name, user,Email)
}
```
The `json.NewDecoder` reads a continuous stream of data (like an open HTTP response body) and decodes the incoming JSON directly into a Go struct using `.Decode(&var)`. When you make an HTTP request, the response body doesn't arrive instantly as a complete string, it streams over the network connection. `json.NewDecoder` reads directly from that stream (`io.Reader`) without needing to load the entire JSON response into memory first using `io.ReadAll()`.

Just like with `sync.Mutex` or `Unmarshal`, `.Decode()` needs to modify the destination variable in-place. If you don't pass a pointer (e.g., `decoder.Decode(issue)` instead of `decoder.Decode(&issue)`), Go will return a runtime error: `json: Unmarshal(non-pointer ...)`.

| Feature | `json.NewDecoder(r).Decode(&v)` | `json.Unmarshal(data, &v)` |
| :--- | :--- | :--- |
| **Input Type** | An `io.Reader` stream (like `resp.Body` or a File). | A `[]byte` slice in memory. |
| **Primary Use Case** | Reading HTTP response/request bodies or large files. | Working with in-memory JSON strings or data already in a byte array. |
| **Performance** | Memory-efficient (streams data directly without buffering all of it). | Requires loading the whole payload into RAM first. |

The `json.Unmarshal` method is used to decode data that is already in a `[]byte` format. It is better used for small programs that is already in memory.
```go
package main

import (
	"encoding/json"
	"fmt"
)

type Issue struct {
	Title string `json:"title"`
	Number int   `json:"number"`
}

func main() {
	jsonBytes := []byte(`{"title": "Fix bug in auth", "number": 42}`)

	var issue Issue
	// Unmarshal requires a pointer (&issue) so it can modify the struct
	err := json.Unmarshal(jsonBytes, &issue)
	if err != nil {
		fmt.Println("Unmarshal error:", err)
		return
	}

	fmt.Printf("Parsed: %+v\n", issue) // Parsed: {Title:Fix bug in auth Number:42}
}
```

The `json.Marshal` function converts a Go struct into a slice of bytes representing JSON data.
```go
type Board struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	TeamId   int    `json:"team"`
	TeamName string `json:"team_name"`
}

board := Board{
	Id:       1,
	Name:     "API",
	TeamId:   9001,
	TeamName: "Backend",
}

data, err := json.Marshal(board)
if err != nil {
	log.Fatal(err)
}
fmt.Println(string(data))
// {"id":1,"name":"API","team":9001,"team_name":"Backend"}
```
#### XML
__Extensible Markup Language__ is a text format for repsenting structured information similar to _HTML_:
```xml
<root>
  <id>1</id>
  <genre>Action</genre>
  <title>Iron Man</title>
  <director>Jon Favreau</director>
</root>
```
It is used for the same things that JSON is used for. However, JSON is more lightweight and arguable easier to read.

### DNS

### URIs