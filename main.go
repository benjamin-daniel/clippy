package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	// hash "github.com/benjamin-daniel/clippy/hash"
	"github.com/benjamin-daniel/clippy/store"
)

func main() {
	// index, _ := hash.GetHash(`should return an hexadecimal`)

	// Create a channel to talk with the OS
	var sigChan = make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	go copyOverTime()

	// Wait for an event
	<-sigChan
	fmt.Print("\n Service Shutting Down\n")

	// time.AfterFunc()

	// store.GetLast()
}

func copyOverTime() {
	for {
		store.AddIfNotPresent()
		time.Sleep(500 * time.Millisecond)
		// time.AfterFunc(500, store.AddIfNotPresent)
	}
}
