package main

import (
	"log"
	"os"

	"github.com/PoohL0ve/lehGO/workspace/gator/internal/config"
)

type state struct {
	cfg *config.Config
}

func main() {
	storedData, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	currentState := &state{
		cfg: &storedData,
	}

	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)

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
