package main

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	fmt.Printf("Please enter your password\n")
	var str string
	fmt.Scanf("%s", &str)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(str), 5)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(hashedPassword))
}
