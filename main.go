package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/benjamin-daniel/clippy/cmd"
	"github.com/benjamin-daniel/clippy/store"
	"github.com/sevlyar/go-daemon"
)

// Path is the path to the folder that holds our data
var Path string = "/usr/local/clippy"

func init() {
	f, err := os.OpenFile(Path+"/test.db", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
}

// To terminate the daemon use:
//  kill `cat sample.pid`
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Just calling ")
		return
	}
	checkStart := os.Args[1:2][0]
	if checkStart != "start" {
		cmd.Execute()
		return
	}
	err := os.MkdirAll(Path, 0777)
	if err != nil {
		log.Fatalf("MkdirAll %q: %s", Path, err)
	}
	cntxt := &daemon.Context{
		PidFileName: "clippy.pid",
		PidFilePerm: 0644,
		LogFileName: "clippy.log",
		LogFilePerm: 0640,
		// WorkDir:     "./",
		WorkDir: Path,
		Umask:   027,
		// Args:        []string{"[go-daemon clippy]"},
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
