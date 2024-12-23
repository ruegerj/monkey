package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/ruegerj/monkey/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! Monkey greets you :)\n", user.Username)
	fmt.Println("Feel free to write some code:")
	repl.Start(os.Stdin, os.Stdout)
}
