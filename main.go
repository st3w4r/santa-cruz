package main

import (
	"context"
	"log"

	"github.com/st3w4r/santa-cruz/cmd"
	"github.com/st3w4r/santa-cruz/config"
)

func main() {

	_, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	cmd.NewCLI().ExecuteContext(context.Background())
}
