package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	argsWithProg := os.Args
	if len(argsWithProg) != 2 {
		fmt.Println("Usage: pwd_tool password")
		return
	}

	passwordArg := strings.TrimSpace(argsWithProg[1])
	passwordBytes := []byte(passwordArg)

	hashedPassword, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(hashedPassword))
}
