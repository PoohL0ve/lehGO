# Learning Plan

## Phase 1: Familiarisation
Complete every exercise on [A Tour of Go](https://go.dev/tour/list). Learn and master the Go Toolchain to be comfortable initialising modules and running programs:
- `go run`
- `go build`
- `go test`
- `go mod`

Finally, read [Effective Go](https://go.dev/doc/effective_go) as a style guide to understand why Go works and is written the way it does. Supplement it with the official Go Code Review Comments document on GitHub (which reflects modern style standards) and explicitly add `go vet` and `golangci-lint` to the toolchain list to catch non-idiomatic code automatically.

## Phase 2: Building Memory with Projects
Try out the exercises on [Gophercises](https://gophercises.com/) before viewing the solutions. Also practice go through [TDD](https://quii.gitbook.io/learn-go-with-tests).

Use `pkg.go.dev` to inspect the source code of the standard library itself. Clicking on functions like `strings.Builder` or `bufio.Scanner` opens the underlying Go implementation, offering a direct look at code written by the Go core team.


## Phase 3: Concurreny and System Thinking
Master goroutines, channels, and system design through rigorous feedback. Use Exercism to get real human feedback and __stop writing PYTHON style code in Go__.

To build on System Thinking, use Coding challenges by John Crickett. Always compile and test concurrent code with the race detector enabled (`go test -race` and `go run -race`). Concurrency bugs are non-deterministic; the race detector makes them explicit. For theory, pair John Crickett's challenges with the book "Concurrency in Go" by Katherine Cox-Buday

## Phase 4: Real World Architecture
Build production-grade services with database integration and deployment.

Enforce a rule of zero external framework dependencies for your first two backend services. Build using only standard library packages (`net/http`, `database/sql`, `context`, `crypto`) before reaching for third-party routers or ORMs like Gin, Fiber, or GORM.

## Phase 5: Domain Specific Paths
To specialize in Go's primary industry domains without relying on step-by-step tutorials, use specification-driven projects. In these challenges, you are given requirements or test suites, but zero code implementation.

1. __Web Backend & Microservices__
    - __Core Skills__: HTTP protocol mechanics, raw SQL query execution, middleware pipelines, JWT/session authentication, gRPC/Protobuf inter-service communication, and graceful shutdowns.
    - __Self-Directed Challenge__: Build an API Gateway from scratch using only net/http and database/sql. It must route requests, enforce token bucket rate-limiting per IP address using Go channels/tickers, handle CORS, and log structured JSON to os.Stdout.
    - Resources (Forced Self-Building):
        - RFC 9110 (HTTP Semantics): Read the actual protocol spec to understand how headers, status codes, and chunked encoding work rather than relying on framework abstractions.
        - Alex Edwards' "Let's Go": Use the table of contents as a feature specification checklist for a backend service, but implement every feature without reading the code snippets.
2. __Systems Programming & Networking__
    - __Core Skills__: Low-level TCP/UDP socket management, binary protocol parsing, memory buffers (bufio, bytes), process signals, and mutex-protected memory state.
    - __Self-Directed Challenge__: Build a Redis-compatible in-memory key-value database. It must accept raw TCP connections, parse the RESP (Redis Serialization Protocol), support GET, SET, and EXPIRE commands, and safely handle concurrent client reads/writes using read-write mutexes (sync.RWMutex).
    - Resources (Forced Self-Building):
        - CodeCrafters.io: Provides automated test suites that run against your local Git repository. It acts as an automated compiler/tester for building systems (Redis, Docker, Git, HTTP server) in Go without providing code solutions.
        - Build Your Own X (GitHub Repository): A curated collection of specifications for rebuilding production software from scratch.
3. __Distributed Systems & Infrastructure__
    - __Core Skills__: Consensus algorithms, leader election, network partitioning, RPC, distributed locking, and event streaming.
    - __Self-Directed Challenge__: Build a distributed fault-tolerant key-value store using a consensus protocol where node failures do not corrupt data or stop writes across the cluster.
    - Resources (Forced Self-Building):
        - MIT 6.824 (Distributed Systems): The entire course curriculum and lab assignments are publicly accessible for free. The labs require implementing Raft consensus and fault-tolerant key-value services in Go using skeleton interfaces and strict automated test suites.
4. __DevOps, Cloud Native & CLI Tooling__
    - __Core Skills__: POSIX flags, stdin/stdout piping, process spawning (os/exec), filesystem traversal (path/filepath), and concurrency worker pools.
    - __Self-Directed Challenge__: Build a high-throughput, concurrent CLI log parser. It must accept multi-gigabyte log files via stdout or file flags, spin up a worker pool scaled to runtime.NumCPU(), parse log lines matching specific regex patterns, and output aggregate error statistics within milliseconds.
    - Resources (Forced Self-Building):
        - Cobra Library Source Code: Study the implementation of [github.com/spf13/cobra](https://github.com/spf13/cobra) to understand how professional CLI tools structure command trees and flags natively in Go.


🚫 Rules to Escape Tutorial Hell
- The 20-Minute Rule: When stuck, struggle for 20 minutes before looking at documentation. The struggle creates the neural pathways for learning.
- Documentation First: pkg.go.dev is your bible. It is often more accurate and up-to-date than any blog post.
- Build Ugly Code First: Do not aim for perfection. Get it working, then refactor. Tutorials often show perfect code instantly, which is unrealistic.
- No "Todo Apps": Build tools you need. If you need a log parser, build it in Go. Personal utility drives deeper learning than generic exercises. 
