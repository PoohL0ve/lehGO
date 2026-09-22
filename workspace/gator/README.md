# Gator

A CLI tool that allows users to:
- Collect **RSS** feeds from across the internet.
- Store collected data in **PostgreSQL**.
- Follow and unfollow RSS feeds added by users.
- View summaries of aggregated posts with links to full posts.

An **RSS feed (Really Simple Syndication)** is a standardized format (XML) for websites to share content updates. Instead of manually visiting sites, a single tool pulls all updates into one place.

**Key Benefits:**
- **Convenience:** No need to manually check multiple websites for new content.
- **Efficiency:** View headlines and snippets without page layouts or ads.
- **Privacy:** Read updates without site-level tracking.
- **Control:** Subscribe exclusively to preferred sources.

---

## Tools

| Tool | Purpose |
| :--- | :--- |
| **Go** | Primary programming language |
| **PostgreSQL** | Relational database storage |
| **Goose** | Database migration management |
| **SQLC** | Type-safe SQL code generator for Go |
| **golangci-lint** | Static analysis and linter |

---

## Config

*Details on state file management and local environment setup.*

State is persisted locally in `~/.gatorconfig.json` using the `internal/config` package.

- **`Config` Struct:** Maps `db_url` and `current_user_name` JSON fields.
- **`Read()`:** Resolves the home directory via `os.UserHomeDir()` and decodes the configuration file.
- **`SetUser()`:** Updates the logged-in user state and saves the formatted JSON to disk with `0600` file permissions.

## Database
__PostgreSQL__ runs locally on macOS using the native system user as superuser (no extra passwords required).
```bash
# Start Postgres server
brew services start postgres@version
```
```bash
createdb gator
## Connect to gator db
psql -d gator
exit # exits psql
```

Useful `psql` commands
|Command      |Action                |
|:-----------:|:---------------------|
|`\c` | Connect / switch to the database|
|`\l` | List all available databases |
|`\dt` | List all current tables in the database |
|`\q` | Quit/exit session |

### Database Migrations
Schema changes are tracked as version-controlled SQL files inside `sql/schema/` using **Goose**.

* **Why Migrations Matter:** They allow systematic, reversible alterations to the database schema (`Up` to apply changes, `Down` to revert) so every developer and deployment environment remains in lockstep.

### Code Generation with SQLC

`sqlc` compiles raw SQL schema files and queries into type-safe Go code inside `internal/database`.

* **Config:** Managed via `sqlc.yaml` in the project root.
* **Schema Source:** `sql/schema/`
* **Query Source:** `sql/queries/`
* **Output Path:** `internal/database/`

### SQL Queries & Code Generation (`sqlc`)

SQL statements live in `sql/queries/` and are compiled into type-safe Go code in `internal/database/`.

#### Generating Go Code
Run from the project root:
```bash
# ensure the uuid is in the mod file
go get github.com/google/uuid
sqlc generate
```
```sql
-- name: GetUser :one
SELECT * FROM users
WHERE name = $1;
```
- `-- name: GetUser :one`: Tells sqlc to generate a method called GetUser returning a single row.

#### Query Reference (`sql/queries/users.sql`)
* **`CreateUser`**: `-- name: CreateUser :one` — Inserts a new user and returns the created record.
* **`GetUser`**: `-- name: GetUser :one` — Selects and returns a single user matching the provided name string (`WHERE name = $1`).
* **`GetUsers`**: `-- name: GetUsers :many` — Returns all user records from the database.
* **`ResetUsers`**: `-- name: ResetUsers :exec` — Deletes all records from the `users` table.

#### Query Reference (`sql/queries/feeds.sql`)
* **`CreateFeed`**: `-- name: CreateFeed :one` — Inserts a new feed associated with a user ID and returns the created record.
* **`GetFeeds`**: `-- name: GetFeeds :many` — Returns all feeds joined with the creator's username (`feeds JOIN users`).
* **`GetFeedByURL`**: `-- name: GetFeedByURL :one` — Looks up a feed record by its URL.

#### Query Reference (`sql/queries/feed_follows.sql`)
* **`CreateFeedFollow`**: `-- name: CreateFeedFollow :one` — Inserts a feed follow record using a CTE and returns the full row along with the linked `user_name` and `feed_name`.
* **`GetFeedFollowsForUser`**: `-- name: GetFeedFollowsForUser :many` — Returns all feed follow rows for a given user ID, joining feed and user names.

### PostgreSQL Driver (`github.com/lib/pq`)

* **Blank Import (`_ "github.com/lib/pq"`):** Used exclusively for its side effect of executing the driver's internal `init()` function to register PostgreSQL support with Go's `database/sql` package.

### State Management

The `state` struct holds application context shared across all command handlers:

```go
type state struct {
    db  *database.Queries
    cfg *config.Config
}
```

#### `feeds` Table
* **`id`**: UUID primary key.
* **`created_at` / `updated_at`**: Creation and modification timestamps.
* **`name`**: Feed name string.
* **`url`**: Unique feed URL.
* **`user_id`**: Foreign key pointing to `users(id)` with `ON DELETE CASCADE`.

## RSS

*Parsing logic for XML feeds and HTTP fetching.*
### RSS Parsing

Example RSS XML Schema:

```xml
<rss xmlns:atom="[http://www.w3.org/2005/Atom](http://www.w3.org/2005/Atom)" version="2.0">
<channel>
  <title>RSS Feed Example</title>
  <link>[https://www.example.com](https://www.example.com)</link>
  <description>This is an example RSS feed</description>
  <item>
    <title>First Article</title>
    <link>[https://www.example.com/article1](https://www.example.com/article1)</link>
    <description>This is the content of the first article.</description>
    <pubDate>Mon, 06 Sep 2021 12:00:00 GMT</pubDate>
  </item>
</channel>
</rss>
```

## Feed Following System

The feed following system introduces a many-to-many relationship between users and RSS feeds using a `feed_follows` join table.

* **Many-to-Many Architecture:** Multiple users can follow the same unique feed URL, and a single user can follow multiple feeds.
* **Cascade Deletions:** Deleting a user or feed automatically purges all corresponding `feed_follows` entries (`ON DELETE CASCADE`).
* **Uniqueness Guarantee:** A composite unique constraint on `(user_id, feed_id)` prevents users from following the same feed multiple times.

```bash
goose -dir sql/schema postgres "postgres://@localhost:5432/gator?sslmode=disable" up
```

## Aggregator

*Background worker loop and feed processing.*

## Commands
The application uses a hand-rolled CLI router mapping command strings to handler functions.

| Command | Usage | Description |
| :--- | :--- | :--- |
| **`login`** | `gator login <username>` | Validates user existence in PostgreSQL and sets the active user in `~/.gatorconfig.json`. Errors if user is not found. |
| **`register`** | `gator register <name>` | Creates a new user in PostgreSQL and sets them as the active user. |
| **`reset`** | `gator reset` | Clears all users from the database for development resetting. |
| **`users`** | `gator users` | Displays all registered users in the database, marking the currently active session. |
|**`agg`** | `gator agg` | Fetches, parses, and outputs RSS feeds.|
| **`addfeed`** | `gator addfeed <name> <url>` | Creates a new feed record linked to the currently logged-in user. |
| **`feeds`** | `gator feeds` | Lists all feeds in the database along with their URLs and creator usernames. |
| **`follow`** | `gator follow <url>` | Follows an existing RSS feed for the current user. |
| **`following`** | `gator following` | Lists all RSS feeds currently followed by the logged-in user. |
| **`unfollow`** | `gator unfollow <url>` | Unfollows an RSS feed for the currently logged-in user. |