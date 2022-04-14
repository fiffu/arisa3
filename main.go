package main

import (
	"arisa3/app"
	"flag"
	"log"
	"os"
)

// Bot parameters
var (
	ConfigFilePath = flag.String("config-file", "", "Config file path")
)

func assertFlags() {
	if *ConfigFilePath == "" {
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	flag.Parse()
	assertFlags()

	err := app.Main(*ConfigFilePath)
	if err != nil {
		log.Fatalf("Failed to start bot, err=%v", err)
	}
	log.Println("Gracefully exit")
}
