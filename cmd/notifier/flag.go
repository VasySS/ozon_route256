package main

import (
	"flag"
	"workshop-1/config"
)

func init() {
	flag.StringVar(&config.KafkaBrokers, "bootstrap-server", "localhost:9092", "kafka broker host and port")

	flag.Parse()
}
