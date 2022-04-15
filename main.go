package main

import (
	"arisa3/app"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/debug"
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

func panicWatch(run func() error) (err, errPanic error) {
	defer func() {
		if r := recover(); r != nil {
			errPanic = fmt.Errorf("panic -> %+v\n%s", r, string(debug.Stack()))
		}
	}()
	err = run()
	return
}

func main() {
	run := func() error {
		flag.Parse()
		assertFlags()

		return app.Main(*ConfigFilePath)
	}
	err, errPanic := panicWatch(run)
	if err != nil {
		log.Fatalf("Failed to start bot, err=%v", err)
	}
	if errPanic != nil {
		log.Fatalf("Runtime error\n%+v", errPanic)
	}
	log.Println("Gracefully exit")
}
