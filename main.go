package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Println("no website provided")
		os.Exit(1)
	}
	if len(args) > 2 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	baseURL := args[1]
	fmt.Printf("starting crawl of: %s...\n", baseURL)

	pages := make(map[string]int)

	crawlPage(baseURL, baseURL, pages)

	fmt.Printf("\n\n")

	for key, val := range pages {
		fmt.Printf("%s : %d\n", key, val)
	}
}
