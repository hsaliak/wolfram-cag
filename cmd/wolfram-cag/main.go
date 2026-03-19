package main

import (
	"fmt"
	"os"

	"wolfapi/pkg/wolframcag"
)

func main() {
	if err := wolframcag.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
