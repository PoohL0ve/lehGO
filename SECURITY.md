# Web Security in Go
This course is mainly taken from boot.dev to understand how to protect web applications from _injection vulnerabilities_ and _cross-site scripting_, as well as authenticating users and preventing data leaks. 

- [Owasp](#owasp)
- [Authentication](#authentication)
- [Authorisation](#authorisation)
- [Untrusted input](#injection)

## OWASP
When building applications, security is **non-negotiable**. Developers make simple mistakes that attackers exploit:
- weak authentication
- insecure defaults
- forgotten debug routes
- outdated dependencies

*Project*: Bearly Secure, Go web app to harden throughout course.

The **Open Worldwide Application Security Project** publishes a new list of common attackers every few years, which can be used as a starting point to securing applications. 

A common example of **A02: Security Misconfiguration** is missing headers in HTTP, such as `X-Content-Type-Options: nosniff`, which informs the browser to respect a response's _Content-Type_ instead of guessing it. For Bearly Secure, this can be done by creating a [middleware](/webplayground/learn-web-security-golang-starter/internal/httpserver/middleware.go).

Web attackers are mainly interested in:
- **Money**: through payments, refunds, or extortion.
- **Data**: emails, names, addresses, credentials to abuse or sell.
- **Access**: control of accounts or admin panels.
- **Disruption**: making services unreliable.

_Note_: They can attack backends through API tools like _cURL_ and Postman.

A **trust boundary** is where parts of a system is separated based on the level of trust. In other words, the backend server is trusted because th developer is in control. All incoming requested should be treated like a potential attack, therefore the server needs to _validate_ the data and _authorise_ the requested operation. __Threat Modeling__ is about examining a system from an attacker's perspective and asking _what could go wrong_. **STRIDE** can be used for 6 common classes of threats:
1. **Spoofing**: pretending to be another user, violating authentication.
2. **Tampering**: changing data at rest or tansit without authorisation, hampering with integrity.
3. **Repudiation**: denying an action because the system lacks reliable records, where an attacker may manipulate logs, leading to an accounting issue.
4. **Information disclosure**: Violates confidentiality, by accessing sensitive data.
5. **Denial of service**: slowing a system or forcing it to be unreliable which violates availability.
6. **Elevation of priveldge**: Tampers with authorisation by gaining permission to resources an attacker should not have.

## Authentication
It is about verifying **WHO** a user is. Applications should not act like email systems which applies case-sensitivity by converting addresses to one format. Instead, apps should be strict about their canonical representation for identities before storing them or looking them up. For Bearly Secure, emails should be lowercase and lack surrounding whitespace. This is done by modifying [`accounts.NormalizeEmail`](./webplayground/learn-web-security-golang-starter/internal/accounts/store.go).

Requests generally use *stateless tokens* or *stateful sessions* as proof of authentication:
|Stateless Authentication    | Stateful Athentication   |
|:---------------------------|:-------------------------|
|All requests come with a credential that the server can validate. This occurs:<br>The server issues a signed token to the client after the client logs in.<br>The client sends the signed token with each request.<br>The server accepts or rejects the request based on the token ithout loading a session record.<br>This means that the server do not depend on stores session states. Many of them use bearer tokens encoded as _JWTS_; however, if someone else get a hold of it they can use it.| The server can identifier the client through session records. When a user logs in, the server creates a session record and sends a session identifier to the client. With future requests the client sends the secret identifier, which the server accepts or rejects as it maps it to a valid session record, ususally in a database or in-memory store.|

#### Sessions and Cookies
A session is a temporary stateful interaction between a user and a system that allows a server to maintain data persistence across multiple requests. As HTTP protocols are stateless, sessions assign a unique __session ID__ to users, that the server users to retrieve stored information, such as login, preferences, shoppint cart contents. This ID is typically stored in a __cookie__. A session table stores:
- unique ID, a secret given to the client.
- user ID (authenticated user)
- Session expiry info
- Session status (revoked or not)

The server reads the Client header for the session ID on every requests and rejects sessions that are missing, expired, or revoked. _Bearly Secure_ implements [sessions](/webplayground/learn-web-security-golang-starter/internal/account/handler.go) to protect user account pages.

A Session ID is just a random, secret string that lives in the server's database to identify a logged-in user. A Cookie is not a session; a cookie is simply the storage mechanism/envelope that the web browser uses to carry that session ID back and forth. A cookie contains elements like a name, value, and optional attributes including _Secure_ (sent only over HTTPS), _HttpOnly_, and _Samesite_ (restricts XSS). HTTP Communication flow:
1. __The Server issues the tag__: When a user logs in successfully, the server replies with an instruction telling the browser to save a cookie:
```HTTP
HTTP/1.1/ 200 OK
Set-Cookie: bootdev_session_id=abc123xyz; Secure; HttpOnly
```
2. __The Browser Stores it__: The broswers places the session id and its contents in the cookie.
3. __The Browser uses it__: when a user later clicks a link to see their account page, the broswer looks up the cookie and attaches it to the request header:
```HTTP
GET /account HTTP/1.1
Host: localhost:3030
Cookie: bootdev_session_id=abc123xyz
```
4. __Server verification__: The backend code extracts the `abc123xyz` and looks it up in the database session table and confirms it has not expired.

For _Bearly Secure_ ensure that the session and cookie records agree by updating the [session_id](./webplayground/learn-web-security-golang-starter/internal/auth/sessions/sessions.go).

Although cookies are convenient for browsers, session cookies are exploitable for attackers, where they can steal a bearer session cookie and reuse it and the server won't know the difference. With cookie flags these attacks can be limited:
- `HttpOnly`: a boolean that prevents JavaScript from reading a cookie through `document.cookie`.
- `Secure`: a boolean that is applied if a cookie is sent using plain HTTP and not HTTPS so no one can read the session_id from the header.
- `SameSite`: can use three different values which restricts when the browser sends cookies with cross-site requests to protect against Cross-Site Request Forgery (CSRF) attacks:
    - `Strict`: The cookie is never sent on cross-site requests, meaning if a user clicks an external link to your site, they will appear logged out until they navigate internally.
    - `Lax`: The cookie is sent only when a user navigates to your site directly via a top-level link (like clicking an ordinary hyperlink), but blocked on background requests like embedded images or iframes.
    - `None`: The cookie is sent on all cross-site requests, but it requires the `Secure` flag so it only travels over encrypted HTTPS connections.
- _Do not make session ids predictable or sequential, as that makes them easy tog guess_.

> [!Important]
> 
>Cookie attributes should always be set.<br>
> Do not rely on frameworks.

```go
// Setting Secure Cookies
package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
)

// Generate a cryptographically secure random session ID
func generateSessionID() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Authenticate user credentials here...
	
	// 2. Generate secure session identifier
	sessionID, err := generateSessionID()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 3. Set a secure session cookie
	cookie := http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		MaxAge:   86400,                          // 1 day in seconds
		HttpOnly: true,                           // Prevents JavaScript access (mitigates XSS)
		Secure:   true,                           // Requires HTTPS (transmits only over encrypted channels)
		SameSite: http.SameSiteStrictMode,        // Blocks cookie from being sent on cross-site requests (mitigates CSRF)
	}

	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Logged in successfully"}`))
}

```

Update the [session_id](./webplayground/learn-web-security-golang-starter/internal/auth/sessions/sessions.go) to include the attributes. Look at the browser developer tools under the network tab, search fro Headers, and ensure the updates are visible for the POST request.

Seesion should have a lifetime, normally in session_timeouts, revocations, and expiration dates. A server should enforce and expiration through the `expires_at` field. Sometimes expiration is not enough, thus revoking a session, which is **invalidating it immediately** is necessary. As Bearly Secure clears the broswers cookie when a user logs out, which is not enough, as the server-side record is still active. To prevent attackers from using the session_id, the current user [session](./webplayground/learn-web-security-golang-starter/internal/accounts/store.go) needs to be revoked when logged out.

#### Password Resets
Passwords are often forgotten. For a password requests, the server generates a random token, embeds in a link and emails it to a user. When the user opens the link the server validates the token and allows the user to set a new password. For security reasons password reset tokens should be unpredictable, short-lived, single-use, bound to one user account, and stored as a hash. This can be done for Bearly secure by ensuring the [generated token](./webplayground/learn-web-security-golang-starter/internal/auth/passwordreset/store.go) is properly secured.

__Learn: SAML and OIDC__.

### More Authentication
__Reauthentication__ needs more proof from a logged in user before they can modify sensitive information like transferring money or changing their password. __Step-up authentication__ needs a stronger factor. Allow Bearly Secure to let users change their emails by entering their current [password](./webplayground/learn-web-security-golang-starter/internal/account/handler.go).

Use password generators opposed to traditional caps, letter, symbols, and numbers passwords. Passowrds should only be changed when there is some type of security risk not often. Security questions like *What is your first pet name?* are ridiculous. Instead use password lenghts and block known passwords.

### Multi-Factor Authentication
__MFA__ requires proof from two or more distinct factor types:
1. *Something you know*: A password or PIN.
2. *Something you have*: A device containing an authenticator secret.
3. *Something you are*: A biometric characteristic.

Mixing them up is ideal, though the last one is not an authenticator within itself. A password and A PIN is silly as they are the same factor. MFAs are useful because an atacker will still need another factor to enter a system. Common second factors are:
- **TOTP** codes: short-lived codes generated by an authenticator app.
- **Out-of-band authentication**: separate device exchanges short-lived secret.
- **Biometric-activated authentication**: device checks fingerprint or face before using stored authenticated secret.
- **SMS codes**: short-lived codes sent via text. Not fully ideal as phone numbers can be hacked.

**Time-Based One-Time Password (TOTP)** generates short-lived codes and the current time usually around 6 characters and a 30 second window. It does not require a separate hardware security key and it can work offline. TOTP works by being built on top of a **shared key** between an app and an authenticator, where both parties compute the expected code for the current time. 
```go
timeStep := (unixTime - t0) / 30
serverComputedCode := hotp(sharedSecret, timeStep)

if serverComputedCode == userSuppliedCode {
  // TOTP code accepted
} else {
  // TOTP code rejected
}
```
__Prevent TOTP-code__ [replay](./webplayground/learn-web-security-golang-starter/internal/auth/mfa/store.go).

A __hardware token__ is a physical authenticator that proves possession of a secret stored on the device. While some display one-time codes **Fast Indentity Online (FIDO)**, a type of hardware token that works with *WebAuthn*, connects to a device and use a public-key cryptography to sign a challenge. Unlike TOTP, it uses *asymmetric cryptography*, where the server never sees the private key, making it phishing resistant.

A __passkey__ is a FIDO credential created to replace passwords. __Device-bound passkeys__ are permanently locked to a single physical hardware element. They cannot be copied or exported. **Web Authentication API (WebAuthn)** is the core browser standard that lets the server talk to a user’s authenticator using public-key cryptography. When a user registers a device-bound passkey, their device generates a unique cryptographic key pair. The private key is securely locked inside the hardware authenticator (and never leaves it), while the public key and a unique credential ID are sent to the Go backend and saved in the database. 
```bash
go get github.com/go-webauthn/webauthn
```

**Recovery** for accounts should be rare, rate limited, and auditable, and the user should be notified through an independent channel after every recovery. Emails can be weak for high-stakes accounts, which attackers know, therefore recovery has to balance safety and usability (don't permanantely lock out legitimate users). Alternatives to email recovery include:
- **Saved recovery codes**: random secrets generated during enrollment and stored locally by the user.
- **Recovery contacts**: prearranged contacts, like an admin who will participate in recovery.
- **Repeated identity proofing**: Repeating enough of the original identity checks to establish that the claimant matches the identity tied to the account. 

## Authorisation
**Authorization** is the security process that determines whether an authenticated user has the necessary permissions to access a specific resource or perform a specific action. While authentication verifies *who you are*, authorization decides *what you are allowed to do*. Once your server knows your identity, authorization checks your roles and rules to permit or deny access.

The **Principle of Least Privilege** dictates that every user, process, program, or system component must operate using only the bare minimum set of permissions necessary to complete its specific task and nothing more. *Over-privilege* occurs when users, services, or API clients accumulate more permissions than they actually need for their daily functions, often remaining unchecked over time. It is one of the leading root causes of severe data breaches. When an over-privileged service token leaks, attackers gain lateral movement across internal databases and infrastructure. Common sources of *over-privilege* include:
- **The *Admin By Default* AntiPattern**: Assigning administrator or root-level permissions to new user accounts or internal microservice tokens during initial setup to "prevent permission bugs," with the intention of restricting them later (which rarely happens)
- **Permission Creep (Role Accumulation)**: As developers or employees change roles or projects over time, new permissions are piled onto their existing accounts without revoking old, unused ones.
- **Overly Broad API Scopes**: Creating lazy, catch-all access tokens (e.g., granting a read-only reporting script full read/write/delete database permissions just because it was easier to configure).
- **Static Hardcoded Secrets & Shared Service Accounts**: Using a single database credential or API key shared across multiple microservices, allowing a compromise in one minor service to grant access to everything else.

Permissions should be *explicit, minimal, and reviewed regularly*. For Bearly Secure [restrict](./webplayground/learn-web-security-golang-starter/internal/admin/handler.go) product editing to administrators.

**Broken Access Control** happens when an application fails to properly restrict what authenticated users can do or see. Instead of verifying permissions on every single sensitive action, the system blindly trusts user-supplied input (like changing an ID in a URL to view someone else's data, known as **IDOR** or Insecure Direct Object Reference). The *rule of thumb* is that every protected action needs a server-side permission check. *Enforece* order [ownership](./webplayground/learn-web-security-golang-starter/internal/orders/handler.go) on ther server in Bearly Secure.

*N.B*: Never trust the client. Always do calculations on the server side.

**Role-Based Access Control (RBAC)** assigns permissions to roles such as developer, admin, and manager, and users are assigned to those roles. *Centralizing key policies* means extracting these authorization rules out of individual API handlers and placing them into a single, authoritative policy engine or middleware layer. Scattering authorization checks (`if user.Role == "admin"`) across hundreds of code files leads to missed checks and drift. Centralizing them ensures a single source of truth for who can do what.
```go
package main

import (
	"net/http"
)

// Centralized role policy map (Role -> Allowed Actions/Permissions)
var rolePolicies = map[string][]string{
	"admin":  {"read", "write", "delete", "manage_users"},
	"editor": {"read", "write"},
	"viewer": {"read"},
}

// CheckPermission verifies if a role possesses a specific permission
func hasPermission(role string, requiredPermission string) bool {
	permissions, exists := rolePolicies[role]
	if !exists {
		return false
	}
	for _, p := range permissions {
		if p == requiredPermission {
			return true
		}
	}
	return false
}

// RequirePermission is a centralized middleware enforcing RBAC
func RequirePermission(requiredPermission string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Assume user role is injected into context by auth middleware
		userRole, ok := r.Context().Value("role").(string)
		if !ok || !hasPermission(userRole, requiredPermission) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
```

**Attribute-Based Access Control (ABAC)** evaluates a dynamic combination of attributes before granting access. These include User attributes (e.g., department, clearance level), Resource attributes (e.g., document classification, owner ID), Action attributes (e.g., read, delete), and Environment attributes (e.g., current time, IP address, device trust state). For Bearly Secure enforce reviews [ownership](./webplayground/learn-web-security-golang-starter/internal/reviews/handler.go) across every review handler.

| Feature / Dimension | Role-Based Access Control (RBAC) | Attribute-Based Access Control (ABAC) |
| :--- | :--- | :--- |
| **Primary Basis** | User's assigned organizational role (e.g., Admin, Viewer) | Dynamic attributes (User, Resource, Environment, Action) |
| **Granularity** | Coarse-grained (broad permissions tied to static groups) | Fine-grained (deeply contextualized, real-time rules) |
| **Complexity** | Simple to implement, manage, and audit for standard apps | Complex to design; often requires policy engines (e.g., OPA/Rego) |
| **Flexibility** | Rigid; requires creating new roles for slight permission tweaks | Highly flexible; adapts to new conditions without changing roles |
| **Best Used For** | Standard SaaS apps with clear, hierarchical user tiers | Enterprise systems requiring strict contextual compliance (healthcare, finance) |

A **Signed URL** is a standard web link that contains temporary, cryptographically secured permissions to access a specific resource (such as a private file, download link, or video stream). It includes a built-in expiration timestamp and a cryptographic signature (usually HMAC) generated using a secret key known only to the server. Sometimes you need to give a user access to a private resource without forcing them through a complex login session or exposing your database credentials. For example, when a user clicks "Download Receipt," you want to grant them access to that exact file for 5 minutes, after which the link automatically becomes invalid. It typically contains 3 main query parameters:
1. **target resource**: identifies the exact resource or action, `/files/1/signed-download`.
2. **expiration timestamp**: a Unix epoch timestamp dictating when the link dies, `expire:1300`.
3. **cryptographic signature**: a hash string proving the URL was generated by your server and hasn't been tampered with. It is a *Hash-based Message Authentication Code (HMAC)* created with `crypto/hmac`.

```go
package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

var serverSecret = []byte("super-secret-server-key-change-in-production")

// 1. Create a signed URL valid for a given duration
func GenerateSignedURL(baseURL string, filename string, duration time.Duration) (string, error) {
	expiresAt := time.Now().Add(duration).Unix()
	
	// Construct the data string to sign: path + expiration
	payload := fmt.Sprintf("file=%s&expires=%d", filename, expiresAt)

	// Generate HMAC-SHA256 signature
	mac := hmac.New(sha256.New, serverSecret)
	mac.Write([]byte(payload))
	signature := hex.EncodeToString(mac.Sum(nil))

	// Build the final URL
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	q := u.Query()
	q.Set("file", filename)
	q.Set("expires", strconv.FormatInt(expiresAt, 10))
	q.Set("sig", signature)
	u.RawQuery = q.Encode()

	return u.String(), nil
}

// 2. Validate the incoming signed URL request
func ValidateSignedURL(filename string, expiresStr string, clientSig string) bool {
	// Parse expiration time
	expiresAt, err := strconv.ParseInt(expiresStr, 10, 64)
	if err != nil {
		return false
	}

	// Check if the link has expired
	if time.Now().Unix() > expiresAt {
		return false // Expired link
	}

	// Recreate the payload string and compute expected HMAC
	payload := fmt.Sprintf("file=%s&expires=%d", filename, expiresAt)
	mac := hmac.New(sha256.New, serverSecret)
	mac.Write([]byte(payload))
	expectedSig := hex.EncodeToString(mac.Sum(nil))

	// Cryptographically compare signatures securely (constant-time comparison)
	return hmac.Equal([]byte(clientSig), []byte(expectedSig))
}

func main() {
	// Generate a link valid for 2 minutes
	signedLink, _ := GenerateSignedURL("https://api.mycompany.com/download", "invoice_99.pdf", 2*time.Minute)
	fmt.Println("Generated Signed URL:\n", signedLink)

	// Simulating validation of the incoming parameters from a request
	isValid := ValidateSignedURL("invoice_99.pdf", "1782349200", "fake_signature_here")
	fmt.Println("\nIs signature valid?", isValid)
}
```

## Injection
It is a category of vulnerabilities where untrusted input is treated as instructions. it occurs where user input is executed by a shell, db, or template engine instead of being treated as data. Injection bugs typically start with *string concatenation*. Parameterisation is the best defence against SQL injection. This practices passes user input as distinct types arguments like `$1` rather than being concatenated to a string, ensuring the database engine parses the query separate from the data preventing SQL injection attacks and allowing query plan caching. In Go, *sqlc* automatically applies paramerisation. It parses raw sql files and generates type-safe Go functions that use prepared statements under the hood. 

*Bearly Secure can now use [prepared](./webplayground/learn-web-security-golang-starter/internal/storefront/store.go) statements.

Injection can be beyon SQL injection, as any input that can be interpreted as instructions is an injection risk, including:
- **Shell commands, template engines, search filters and expressions**

**Validation** and **sanitisation** solves two different problems:
- *Validation*: occurs at the application boundary when input first arrives where it rejects invalid input but asking *is this input allowed*?
- *Sanitisation*: transforms input based on a defined policy like removing extra spaces.

Bearly Secure can be updated to accept only [canonical](./webplayground/learn-web-security-golang-starter/internal/cart/handler.go) whole number strings. Also [trim](./webplayground/learn-web-security-golang-starter/internal/reviews/handler.go) reviews for the app.

**Zip Slip** is a *file-overwrite vulnerability* where a malicious compressed archive (like a `.zip` file) contains file names with directory traversal sequences (such as `../`). When an application extracts the archive without validating the file paths, the extracted files *"slip"* out of the intended destination folder and overwrite critical files elsewhere on the file system. It was discovered by [Synk](https://security.snyk.io/research/zip-slip-vulnerability) in 2018. *f an attacker can upload a specially crafted zip file that overwrites a system binary, a configuration file, or an application script (like a cron job or source code file), they can achieve **Remote Code Execution (RCE)** or take complete control of the server.*
```go
// DANGEROUS: If f.Name is "../../secret.txt", it escapes the target directory
destPath := filepath.Join(targetDir, f.Name)
```

Archive extraction also needs *resource limits*, otherwise it can lead to **zip bomb** attacks, the ability to drastically destroy the performance of a server by exhausting its disk space or memory as a type of *denial-of-service-attack*. 

For Bearly secure, [reject](./webplayground/learn-web-security-golang-starter/internal/uploads/archive.go) archives when an entry could escape or redirect extraction.

**LLM Prompt Injection** happens when untrusted user input tricks an AI model into overriding its original system instructions and executing malicious commands instead. In traditional programming, code (logic) and data (input) are stored separately. Large Language Models process everything—system prompts, developer guidelines, and user inputs—as a single continuous stream of text. Because the model cannot inherently tell the difference between "rules given by the developer" and "text typed by the user," attackers can inject malicious commands directly into the prompt stream. To prevent prompt injection, you must **strictly isolate system instructions** from user-submitted data using structured API roles rather than raw string concatenation.
```go
// DANGEROUS: If userInput contains "Ignore previous rules...", the LLM reads it as a command
prompt := "Summarize this text: " + userInput
response := callLLM(prompt)
```
Modern LLM APIs (like OpenAI, Anthropic, or local models via Ollama) use structured JSON payloads with distinct roles. This tells the underlying model parser: "Treat this system block as immutable code, and treat this user block strictly as inert data to process."
```go
// SECURE: Instructions and user data are sent in separate, isolated slots
payload := map[string]any{
    "model": "gpt-4o",
    "messages": []map[string]string{
        {
            "role":    "system",
            "content": "You are a helpful assistant. Never reveal system secrets.",
        },
        {
            "role":    "user",
            "content": userInput, // The LLM treats this purely as data, not instructions
        },
    },
}
```
For Bearly Secure, give trusted instructions and customer input separate message [roles](./webplayground/learn-web-security-golang-starter/internal/assistant/service.go). For limiting LLM Prompt injections, *message roles* are good but they do not outright prevent it. One way to prevent the injections is to *limit tool access*. LLMs are created with the ability to call tools, like unctions that search the web, access files, call APIs, or execute commands.

File uploads should also be [protected](./webplayground/learn-web-security-golang-starter/internal/uploads/documents.go).

## Client-Side Security
**Cross-Site Scipting (XSS)** occurs when user controlled input becomes executable code in the browser. `Httponly` only prevents JS from reading a cookie through `document.cookie` but not other scripting like stopping the injected script from reading page data or making authenticated same-origin requests as the victim. The defense is to encode untrusted values for the context where they're rendered. Also learn Cross-Site Forgery beyond SameSite attribute.

A **Cross-Site Request Forgery (CSRF) Token** is an unpredictable secret generated by the server for a client's session which is typical used in HTTP headers and forms. It esist to ensure that a request originated from an application's UI and not from a malicious external page.
> [!EXAMPLE]
> It's like logging into a website: bank.com where the server
> creates a session_id cookie and generates a random token (CSRF).
>
> The bank.com site would include a hidden input for a transfer funds form:
> `<input type="hidden" name="csrf_token" value="x9K2mP">`
>
> When the form is submitted the server checks whether the token in the form matches the token tied to the session body, and executes if they match.
>
> If an outside site like getdata.com tries to submit a hidden form to bank.com, the browser sends the session_id cokkie but getdata.com cannot read the csrf_token value `x9K2mP`, thus the server rejects the request with `403 Forbidden`.
>
> However, XSS can bypass CSRF tokens by reading one from the DOM and submitting a valid request. You need to prevent both vulnerabilities!

Check that the token [match](./webplayground/learn-web-security-golang-starter/internal/auth/sessions/) for the forms being submitted at Bearly Secure.




