package main

import (
	"log"

	"github.com/arxon31/bank/internal/app"

	"github.com/arxon31/bank/config/prod"
)

func main() {
	cfg, err := prod.New()
	if err != nil {
		log.Fatalf("config error: %s", err)
	}

	app.RunPublisher(cfg)
}
