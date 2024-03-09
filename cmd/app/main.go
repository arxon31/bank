package main

import (
	"log"

	"github.com/arxon31/bank/config"
	"github.com/arxon31/bank/internal/app"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("config error: %s", err)
	}

	app.Run(cfg)
}
