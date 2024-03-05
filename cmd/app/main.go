package main

import (
	"github.com/arxon31/bank/config"
	"github.com/arxon31/bank/internal/app"
	"log"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("confir error: %s", err)
	}

	app.Run(cfg)
}
