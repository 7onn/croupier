package main

import (
	"fmt"
	"os"

	"github.com/7onn/croupier/app/croupier-api/webserver"
	"github.com/7onn/croupier/internal/croupier"
)

func init() {}

func main() {
	hydratedConfig := webserver.Config{
		Croupier: croupier.Config{},
		Port:     5000,
	}

	srv, err := webserver.New(hydratedConfig)

	if err != nil {
		fmt.Printf("Startup error: %s\n", err)
		os.Exit(1)
	}

	srv.Start(BuildPipeline)
}
