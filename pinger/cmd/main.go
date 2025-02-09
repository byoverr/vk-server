package main

import (
	"pinger/internal/config"
	"pinger/internal/container"
	"time"
)

func main() {
	cfg := config.Load()

	for {
		container.ProcessContainers(cfg)
		time.Sleep(time.Duration(cfg.Pinger.Interval) * time.Second)
	}
}
