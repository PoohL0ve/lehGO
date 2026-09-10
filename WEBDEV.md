# Backend Web Development Using Go
__Contents__:
- [__HTTP Clients__](#http-clients)
    - [JSON](#json)
    - [DNS](#dns)
	- [URIs](#uris)
	- [Headers](#methods)
	- [Methods](#methods)
	- [HTTPS](#https)
	- [Errors](#errors)
	- [cURL](#curl)
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
Computers on a network communicate using __numeric IP addresses__ to route data packets across hardware interfaces. Humans are terrible at remembering long strings of numbers. __Domain Names__ are the human readable format of an IP address. The __Domain Name System__ is a distributed, hierarchical database that translates human-readable domain names like boot.dev into machine-readable IP addresses like `198.51.100.1`. DNS allows users to type names while machines route numbers.

Deploying a website on the internet only requires for steps:
1. Create or access a server that hosts the files of the web site and connect it to the internet.
2. Acquire a domain name.
3. Connect the domain name to an IP address of the server.
4. The server is accessible via the internet.

The _domain name_ is just the host name, such `bootdev.com` and not the other parts of the URL. The `net/url` package is Go's standard library tool for parsing, constructing, and manipulating __URLs__ and query parameters safely. URLs are not just simple strings as they have formatting rules governed by internet standards like __RFC 3986__, which the package utilises. This allows URLs to be broken down into finely types structs where they can be safely used. __WHY__: Constructing URLs manually with string concatenation (e.g., `url + "?key=" + val`) easily leads to bugs, broken encodings, or security vulnerabilities like query parameter injection.

__🧠 Mental Model__: The `url.URL` object can be seen as a disassembled address label:
```plaintext
https://  user:pass@  api.example.com :8080  /v1/users  ?sort=desc&limit=10  #profile
└─Scheme─┘ └─Userinfo─┘ └──────Host─────┘└─Port─┘ └──Path───┘ └─────RawQuery─────┘  └─Fragment┘
```
- __Scheme/Protocol__: Defines the rules used to access the resource (http, https, ftp, mailto). Not all URL protocols require the `//` authority components such as `mailto:` which does not use it.
- __Authority__: Contains optional credentials (user:pass@), the domain host (example.com), and an optional port (:8042).
- __Path__: Hierarchical data identifying the specific resource on the server (/over/there).
- __Query__: Optional non-hierarchical parameters (?name=ferret).
- __Fragment__: Optional sub-resource client-side anchor (#nose).
- __Port__: The virtual points managed by a computer's OS where network connections are made. They numbers from _0 - 65, 535_ where port _0_ is reserved for the system's API. Whenever you connect to another computer over a network, you're connecting to a specific port on that computer, which is listened to by a program on that computer. A port can only be used by one program at a time, which is why there are so many possible ports. Most times the default port for HTTP (80) and HTTPS (443) are used on the internet, which aren't really visible.

When a raw string is parsed with `url.Parse()`, Go breaks that single string into a structured `*url.URL` object where each part is an explicit field you can inspect or rewrite. Common methods and functions from the `net/url` package:
| Function / Method | Example Code | Why & When To Use It |
| :--- | :--- | :--- |
| **`url.Parse`** | `u, err := url.Parse("https://api.example.com/search?q=go")` | **Parsing:** Converts a raw URL string into a structured `*url.URL` object. Use this before inspecting or modifying components of a URL. |
| **`u.Query()`** | `params := u.Query()` | **Extracting Params:** Parses the query string (`?a=1&b=2`) into a `url.Values` map (`map[string][]string`) so you can easily read parameters. |
| **`params.Get()`** | `val := params.Get("q")` | **Reading Key:** Retrieves the first value associated with a given key in a `url.Values` map. Returns `""` if the key doesn't exist. |
| **`params.Set()`** | `params.Set("limit", "10")` | **Setting/Overwriting Key:** Sets a key-value pair in `url.Values`, replacing any existing values for that key. |
| **`params.Add()`** | `params.Add("filter", "active")` | **Appending Key:** Appends a new value to a key in `url.Values` without overwriting existing values under that same key. |
| **`params.Encode()`** | `u.RawQuery = params.Encode()` | **URL Encoding:** Converts `url.Values` back into an HTTP-safe, URL-encoded query string (e.g., handles spaces like `%20` or `+`). |
| **`u.String()`** | `fullURL := u.String()` | **Reconstructing:** Reassembles the entire `*url.URL` struct back into a complete, valid URL string ready for an HTTP request. |
| **`url.QueryEscape`** | `escaped := url.QueryEscape("hello world & co")` | **Manual Encoding:** Encodes a standalone string so it can be safely used inside a URL path or query parameter. |

__Example__: (Taken from boot.dev)
```go
parsedURL, err := url.Parse("https://homestarrunner.com/toons")
if err != nil {
	fmt.Println("error parsing url:", err)
	return
}
// Extract the hostname/domain name
parsedURL.Hostname()
```
The [__Internet Corporation for Assigned Names and Numbers (ICANN)__](https://www.icann.org/) is organisation that manages DNS for the entire internet. Each computer internet requests contacts one of ICANN's _root nameservers_ whose address is configured into each computer system. From there, that nameserver can gather the domain records for a specific domain name from their distributed DNS database.

A __subdomain__ allows a domain to route network traffic to different servers and resources. For example, `bea.findher.com` is a subdomain for `findher.com`.

Building a Safe URL in Go:
```go
package main

import (
	"fmt"
	"net/url"
)

func main() {
	// 1. Parse a base URL
	baseURL, err := url.Parse("https://api.boot.dev/v1/courses")
	if err != nil {
		return
	}

	// 2. Extract existing query parameters (returns url.Values map)
	params := baseURL.Query()

	// 3. Add or update parameters safely (handles encoding automatically!)
	params.Set("sort", "name desc") // Space gets properly encoded as + or %20
	params.Add("tag", "golang")
	params.Add("tag", "backend")

	// 4. Encode params back into the URL struct's RawQuery field
	baseURL.RawQuery = params.Encode()

	// 5. Convert full struct back to string
	fmt.Println(baseURL.String())
	// Output: https://api.boot.dev/v1/courses?sort=name+desc&tag=golang&tag=backend
}
```
Easy ways the safety can be broken:
1. __Directly Mutating `Query()`: `u.Query()` returns a copy of the URL's query parameters as a map. Modifying that map directly does not update the URL until `params.Encode()` is assigned back to `u.RawQuery`.
2. __Using String Concatenation for Query Parameters__: Never build queries by concatenating strings like `baseURL + "?search=" + userInput`. If userInput contains spaces, &, or ?, it will break the URL structure or create security bugs. Always use `url.Values` and `.Encode()`

### URIs
A __Uniform Resource Identifier (URI)__ is a string os characters that identifies a resource either by name, location, or both specifically over the internet. To fetch or modify data over a network, servers and clients need a unique, unambiguous address system. URIs define that universal syntax so web systems can pinpoint precisely what resource is being acted upon. There are two common types of URIs:
1. __Uniform Resource Locator (URL)__: The _how_ and _where_ to find something. 
2. __Uniform Resource Name (URN): The _what_ that something is like an isbn number.

__Example Extracting Components of a URI__:
```go
func newParsedURL(urlString string) ParsedURL {
	parsedUrl, err := url.Parse(urlString)
	if err != nil {
		return ParsedURL{}
	}

	var username string 
	var password string 

	if parsedUrl.User != nil {
		username = parsedUrl.User.Username()
		password, _ = parsedUrl.User.Password()
	}

	return ParsedURL{
		protocol: parsedUrl.Scheme,
		username: username,
		password: password,
		hostname: parsedUrl.Hostname(),
		port:     parsedUrl.Port(),
		pathname: parsedUrl.Path,
		search:   parsedUrl.RawQuery,
		hash:     parsedUrl.Fragment,
	}
}
```
This shows that there are typical 8 components to a URL, but not all are required.
![URL Components](images/url_components.png)
- Image BY <a href="https://www.boot.dev/lessons/2d85b1ef-5577-4f65-88ca-e4264059b7af">__Boot.dev__</a>.

### Headers
__HTTP Headers__ are _key-value_ pairs sent in both HTTP requests and responses to pass metadata alongside the main body payload. The body contains the actual message (like a JSON object), but the server and client need instructions on how to handle that message such as what format the data is in, how to authenticate the user, or whether to cache the content. It's the label of a package.

Common Headers Developers Use:
| Header Name | Type | Purpose & Common Values |
| :--- | :--- | :--- |
| **`Content-Type`** | Both | Specifies the media format of the body (e.g., `application/json`, `text/html`, `multipart/form-data`). |
| **`Accept`** | Request | Tells the server what format the client prefers to receive back (e.g., `application/json`). |
| **`Authorization`** | Request | Carries authentication credentials (e.g., `Bearer <token>`, `Basic <base64>`). |
| **`User-Agent`** | Request | Identifies the software making the request (e.g., `Mozilla/5.0`, `PostmanRuntime`, `Go-http-client/1.1`). |
| **`Set-Cookie`** | Response | Asks the browser/client to store a session cookie for future requests. |
| **`Location`** | Response | Tells the client where to redirect when returning a 3xx status code. |

The `net/http` package provides tools like `Header` to interact with HTTP headers. The _Header_ type is a map of string slices `map[string][]string`.

__Example__:
```go
func main() {
	req, err := http.NewRequest("GET", "https://api.boot.dev/v1/courses", nil)
	if err != nil {
		return
	}

	// 1. Setting a header (Overwrites existing values)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer MY_API_KEY")

	// 2. Adding a header (Appends to existing values for multi-value headers)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept", "text/plain")

	// 3. Reading a header (Case-insensitive!)
	authHeader := req.Header.Get("authorization") // Handles canonical formatting automatically!
	fmt.Println("Auth:", authHeader)

	// 4. Deleting a header
	req.Header.Del("Authorization")
}
```
__Common Mistakes__:
1. __Case Sensitivity & CanonicalMIMEHeaderKey__: HTTP specification states headers are case-insensitive (Content-Type vs content-type). Go automatically normalizes header keys using canonical formatting (capitalizing the first letter and any letter after a hyphen).
	- If you read headers using req.Header.Get("content-type"), Go automatically converts it to "Content-Type" for map lookup.
	- The Pitfall: If you access the map directly with map syntax req.Header["content-type"], it will fail and return nil because Go map lookups are strictly case-sensitive! Always use .Get() or .Set().
2. __Setting Headers on `http.Get()`: Shorthand functions like `http.Get()` or `http.Post()` do not allow you to attach custom headers. If you need to send an Authorization or custom User-Agent header, you must use `http.NewRequest()` and `client.Do()`.

The command `CMD + Opt + I` opens the developer tools.

### Methods
__HTTP Methods__ (also called HTTP Verbs) are standardized commands defined in the HTTP specification that tell the server what primary action to perform on a given resource URL. Instead of creating dozens of custom API endpoints like `/getUser`, `/createUser`, or `/deleteUser`, HTTP methods allow you to use a single endpoint (`/users`) and convey the intended operation directly through the request verb: __GET /users, POST /users, DELETE /users/123__. The server inspects the method verb first to determine what action to take with the request.

__Types of HTTP Verbs__:
| Method | CRUD Action | Safe? | Idempotent? | Expects Request Body? | Primary Purpose |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **`GET`** | Read | **Yes** | **Yes** | No | Retrieves data without modifying anything on the server. |
| **`POST`** | Create | No | No | **Yes** | Submits new data to create a new resource on the server. |
| **`PUT`** | Update | No | **Yes** | **Yes** | Replaces an entire existing resource with a new representation. |
| **`PATCH`** | Update | No | No | **Yes** | Applies partial modifications/updates to an existing resource. |
| **`DELETE`**| Delete | No | **Yes** | Optional / Rare | Removes a specified resource from the server. |
| **`HEAD`** | Read Metadata | **Yes** | **Yes** | No | Same as `GET`, but returns ONLY response headers (no body). |

HTTP Methods are governed by two critical terms:
1. __Safe Methods__: Methods that do not modify the state of the server, such as `GET`.
2. __Indempotent Methods__: Methods where making the exact same request once vs. 100 times in a row yields the same server state result.

For `POST` and `GET` requests, `http.Get` and `http.Post` can be used respectively, but for `PUT`, `PATCH`, and `DELETE`, the `http.NewRequest()` is needed:
```go
package main

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type User struct {
	Name string `json:"name"`
}

func main() {
	user := User{Name: "Alice"}
	bodyBytes, _ := json.Marshal(user)

	// 1. Create a PUT request
	req, err := http.NewRequest("PUT", "https://api.example.com/users/123", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// 2. Execute via client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
}
```
Go offers three ways to send a `GET` requests:
1. `http.Get(url)`: The quick and short way to execute a simple request that uses the built-in `httpDefaultClient`. It should be used for quick CLI tools, test scripts, or one-off exploratory scripts where default settings are fine and custom headers aren't needed.
	- __Drawback__: Uses `http.DefaultClient`, which has no timeout. If the destination server hangs, the application hangs forever.
2. `http.NewRequest("GET", ...) + client.Do(req)`: This is the _production standard_ that build an explicit `*http.Request` struct first, allowing you to attach custom headers, auth tokens, or context timeouts before passing it to an `http.Client`. It should be used for production backend microservices, interacting with authenticated APIs, or when you need fine-grained control over request headers.
```go
// 1. Define a client with a strict safety timeout
client := &http.Client{Timeout: 5 * time.Second}

// 2. Build the request object
req, err := http.NewRequest("GET", "https://api.example.com/items", nil)
if err != nil {
    return err
}

// 3. Attach custom metadata / headers
req.Header.Set("Authorization", "Bearer TOKEN_HERE")
req.Header.Set("Accept", "application/json")

// 4. Execute
resp, err := client.Do(req)
if err != nil {
    return err
}
defer resp.Body.Close()
```
3. `http.NewRequestWithContext(ctx, "GET", ...)`: The Concurrent/Cancellation Pattern that binds an HTTP request to a Go `context.Context` object. It should be used for server requests that spawn outgoing API calls, or concurrent applications where you need to cancel an in-flight network request if another routine fails or if a user closes their connection.
```go
// Automatically cancels the GET call if it takes longer than 2 seconds
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()

req, err := http.NewRequestWithContext(ctx, "GET", "https://api.example.com/items", nil)
if err != nil {
    return err
}

resp, err := client.Do(req)
if err != nil {
    return err // Triggers if context timeout expires before server answers
}
defer resp.Body.Close()
```
| Approach | Syntax Pattern | Best Used For | Key Advantage / Tradeoff |
| :--- | :--- | :--- | :--- |
| **`http.Get`** | `http.Get(url)` | Quick scripts, prototyping, simple public reads. | **Pros:** Zero setup.<br>**Cons:** No timeout, no custom headers. |
| **`http.NewRequest` + `Do`** | `req, _ := http.NewRequest("GET", url, nil)`<br>`client.Do(req)` | **Production standard.** APIs requiring auth tokens, custom headers, or global client timeouts. | **Pros:** Full control over headers and client configuration.<br>**Cons:** Requires slightly more boilerplate. |
| **`http.NewRequestWithContext`** | `req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)` | Microservices, web servers, concurrent routines with timeouts/cancellation. | **Pros:** Prevents wasted resources by canceling requests when context expires.<br>**Cons:** Requires managing Go contexts. |

__HTTP Status Codes__ are 3-digit integer responses issued by a server to communicate the outcome of a client's request. The client needs an instant, standardized way to know if the request succeeded, if data was missing, if authentication failed, or if the server crashed—without having to parse the response body text first. They grouped into categories:
- `100-199`: Informational responses. These are very rare.
- `200-299`: Successful responses. Hopefully, most responses are 200's!
- `300-399`: Redirection messages. These are typically invisible because the browser or HTTP client will automatically do the redirect.
- `400-499`: Client errors. You'll see these often, especially when trying to debug a client application
- `500-599`: Server errors. You'll see these sometimes, usually only if there is a bug on the server.

__Common Codes Developers Should Know__:
| Code | Name | Category | Primary Meaning & Backend Use Case |
| :--- | :--- | :--- | :--- |
| **`200`** | **OK** | Success | Standard response for successful `GET`, `PUT`, or `PATCH` requests. |
| **`201`** | **Created** | Success | Request succeeded and a new resource was created (commonly returned after a successful `POST`). |
| **`204`** | **No Content** | Success | Request succeeded, but there is no payload/body to return (common for `DELETE` operations). |
| **`301`** | **Moved Permanently** | Redirection | The requested resource has permanently moved to a new URL (specifies `Location` header). |
| **`302`** | **Found / Temporary Redirect** | Redirection | The resource is temporarily at a different URL. |
| **`304`** | **Not Modified** | Redirection | Tells the client their cached copy is still fresh; no need to re-download the body. |
| **`400`** | **Bad Request** | Client Error | The server cannot parse the request (e.g., malformed JSON payload or missing parameters). |
| **`401`** | **Unauthorized** | Client Error | Authentication is required (missing or invalid API token/credentials). |
| **`403`** | **Forbidden** | Client Error | The client is authenticated, but lacks permissions/rights to access the resource. |
| **`404`** | **Not Found** | Client Error | The requested endpoint or database record does not exist on the server. |
| **`405`** | **Method Not Allowed** | Client Error | The endpoint exists, but doesn't support the HTTP verb used (e.g., sending `POST` to a `GET`-only route). |
| **`409`** | **Conflict** | Client Error | Request conflicts with current server state (e.g., trying to register a duplicate email). |
| **`429`** | **Too Many Requests** | Client Error | The client has exceeded rate limits; stop spamming requests! |
| **`500`** | **Internal Server Error** | Server Error | The server encountered an unhandled exception or bug (e.g., unhandled panic, database crash). |
| **`502`** | **Bad Gateway** | Server Error | An upstream service or reverse proxy (like NGINX) received an invalid response from the backend. |
| **`503`** | **Service Unavailable** | Server Error | The server is temporarily down due to maintenance or extreme traffic overload. |
| **`504`** | **Gateway Timeout** | Server Error | An upstream server took too long to respond to a proxy/gateway. |

Go provides clear constants in the `net/http` package where the `http.Response` struct has a `.StatusCode` property so developers do not have to manually write status codes and their meaning:
```go
package main

import (
	"fmt"
	"net/http"
)

func handleResponse(resp *http.Response) {
	// Instead of checking against 200, use Go's readable constants
	if resp.StatusCode == http.StatusOK {
		fmt.Println("Success!")
	}

	if resp.StatusCode == http.StatusCreated {
		fmt.Println("Resource created!")
	}

	if resp.StatusCode == http.StatusNotFound {
		fmt.Println("Resource missing!")
	}

	if resp.StatusCode >= http.StatusInternalServerError {
		fmt.Println("Server error on their end!")
	}
}
```

### Paths

### HTTPS
### Errors
### cURL