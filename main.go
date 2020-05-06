package main

import (
	"log"
	"time"

	"github.com/benjamin-daniel/clippy/store"
	"github.com/sevlyar/go-daemon"
)

// To terminate the daemon use:
//  kill `cat sample.pid`
func main() {
	cntxt := &daemon.Context{
		PidFileName: "clippy.pid",
		PidFilePerm: 0644,
		LogFileName: "clippy.log",
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
		Args:        []string{"[go-daemon sample]"},
	}

	d, err := cntxt.Reborn()
	if err != nil {
		log.Fatal("Unable to run: ", err)
	}
	if d != nil {
		return
	}
	defer cntxt.Release()

	log.Print("- - - - - - - - - - - - - - -")
	log.Print("daemon started")

	copyOverTime()
}

func copyOverTime() {
	for {
		store.AddIfNotPresent()
		time.Sleep(500 * time.Millisecond)
	}
}
