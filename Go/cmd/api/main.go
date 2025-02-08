package main

import (
	"fmt"
	"log/slog"
	"os"
)

func main() {
	if err := run(); err != nil {
		slog.Error("an error occurred", "error", err)
		os.Exit(1)
	}

	os.Exit(0)
}

func run() error {
	fmt.Println("Hello, World!")

	return nil
}
