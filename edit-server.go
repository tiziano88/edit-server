package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"

	"github.com/tiziano88/completecommand"
)

var (
	port    = flag.Int("port", 8989, "port to listen on.")
	command = flag.String("command", "edit-in-gvim.sh", "command to run to open the editor. It will receive the temporary file name as first argument.")
)

func main() {
	flag.Parse()
	completecommand.Complete()

	http.HandleFunc("/", handler)
	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		log.Fatalf("error binding socket: %v", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	b := make([]byte, 100000)

	n, err := r.Body.Read(b)
	if err != nil {
		log.Print("error reading request body: %v", err)
		return
	}

	f, err := ioutil.TempFile("", "edit-server-")
	if err != nil {
		log.Print("error creating temporary file: %v", err)
		return
	}
	defer f.Close()

	f.Write(b[:n])

	err = exec.Command(*command, f.Name()).Run()
	if err != nil {
		log.Print("error executing edit command %q: %v", *command, err)
		return
	}

	if err := f.Sync(); err != nil {
		log.Print("error reading temporary file %q: %v", f.Name(), err)
		return
	}

	if _, err := f.Seek(0, 0); err != nil {
		log.Print("error reading temporary file %q: %v", f.Name(), err)
		return
	}

	n, err = f.Read(b)
	if err != nil {
		log.Print("error reading temporary file %q: %v", f.Name(), err)
		return
	}

	w.Write(b[:n])
}
