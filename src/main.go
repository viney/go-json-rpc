package main

import (
	"flag"
	"fmt"

	"engine/user"
)

func main() {
	flag.Parse()

	name, email, err := user.Query(flag.Arg(0))
	if err != nil {
		fmt.Println("user.Query: ", err)
		return
	}

	fmt.Println(name, email)
}
