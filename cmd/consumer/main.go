package main

import (
	"log"

	"github.com/arxon31/bank/config/cons"

	"github.com/arxon31/bank/internal/app"
)

func main() {
	cfg, err := cons.New()
	if err != nil {
		log.Fatalf("config error: %s", err)
	}

	app.RunConsumer(cfg)
}
