package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/PoohL0ve/lehGO/workspace/gator/internal/config"
	"github.com/PoohL0ve/lehGO/workspace/gator/internal/database"
	_ "github.com/lib/pq" // Blank import: prevent compiler from throwing unused error
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	storedData, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	db, err := sql.Open("postgres", storedData.DBURL)
	if err != nil {
		log.Fatalf("Error connecting to databse: %w", err)
	}

	dbQueries := database.New(db)

	currentState := &state{
		db:  dbQueries,
		cfg: &storedData,
	}

	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", handlerAddFeed)
	cmds.register("feeds", handlerFeeds)

	if len(os.Args) < 2 {
		log.Fatal("Error: not enough arguments provided")
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	cmd := command{
		Name: cmdName,
		Args: cmdArgs,
	}

	err = cmds.run(currentState, cmd)
	if err != nil {
		log.Fatal(err)
	}
}
