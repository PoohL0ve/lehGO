# Pokedex
- Project from boot.dev

A __REPL__ (read-eval-print loop) or language shell is an interactive CLI program that takes users input and executes it to produce an output. This project is created using the `PokéAPI` to fetch data.

## Create the Package
To begin the project the _pokedexcli_ repo was created along with a _main.go_ file. The `go mod init PACKAGE_NAME` command was then executed to create the `go.mod` file. 

The `tee` command copies to the standard output.
```bash
go run . | tee repl.log
```

## Use TDT For a Function
A `cleanInput()` was created to store given phrases in a slice where each word in the phase was separated by a space using the `Fields` method of the `strings` package. This method was tested using a __Table-driven Test__ to ensure the method works as it should. 

## Create a Command Map Registry
A composite data-structure was created to define the boundaries of commands in the `cliCommand` struct which allows users to type in defined commands for specific purposes. This was coupled with two callback functions that determined the actions the user can take: exit, or help. 

## Utilising State
A new `internal` package was developed to deal with the state of the progrom. In it contains a `pokeapi` folder that has three main files:

__client.go__:
| Component | Purpose |
| :--- | :--- |
| **`package pokeapi`** | Places this file inside your isolated `pokeapi` package. |
| **`type Client struct`** | Wraps `http.Client`. Hides internal network settings from the rest of the app. |
| **`time.Duration`** | Configures connection timeouts (e.g., `5 * time.Second`) so bad network requests fail gracefully. |
| **`NewClient()`** | Exported constructor (capitalized) allowing `package main` to initialize the API client cleanly. |

__location_list.go__ takes optional pagination URLs, executes the HTTP request using the Client's configured httpClient, decodes the JSON response payload, and returns clean Go data to the caller:
- __File Documentation__: Explains that this file handles fetching and unmarshaling location-area endpoint data.
- `LocationAreasResponse Struct`: Defines the JSON response contract (using struct tags).
- `ListLocations Method`: Attached to `c *Client`, accepts a `pageURL *string`, makes the `GET` request, and handles error checks and JSON decoding.