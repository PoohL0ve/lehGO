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

*Database schema, connection strings, and migration details.*

## RSS

*Parsing logic for XML feeds and HTTP fetching.*

## Aggregator

*Background worker loop and feed processing.*

## Commands
The application uses a hand-rolled CLI router mapping command strings to handler functions.

| Command | Usage | Description |
| :--- | :--- | :--- |
| **`login`** | `gator login <username>` | Sets the active user in `~/.gatorconfig.json`. |