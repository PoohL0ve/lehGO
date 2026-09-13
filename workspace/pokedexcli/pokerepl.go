package main

import "strings"

func cleanInput(text string) []string {
	lowerText := strings.ToLower(text)
	// Fields methods splits the string by whitespace
	textSlice := strings.Fields(lowerText)
	return textSlice
}
