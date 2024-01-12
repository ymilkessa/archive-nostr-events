package main

import (
	"fmt"
)

func getUserInput(prompt string) string {
	var input string
	fmt.Print(prompt)
	fmt.Scan(&input)
	return input
}

func GetNip05FromUser() string {
	return getUserInput("Enter a NIP05 identifier: \n>")
}
