package main

import (
	"accountbalance/cmd"
	"fmt"
)

func main() {
	var err error
	err = cmd.Execute()
	if err != nil {
		panic(fmt.Sprintf("cmd.Execute err: %v", err))
	}
}
