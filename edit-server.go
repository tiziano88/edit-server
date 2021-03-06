package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"

	"github.com/tiziano88/completecommand"
)

var (
	port    = flag.Int("port", 7878, "port to listen on.")
	command = flag.String("command", "./edit-in-gvim.sh", "command to run to open the editor. It will receive the temporary file name as first argument.")
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
	log.Printf("request: %#v", r.Header)

	url := r.Header.Get("X-Url")

	b := make([]byte, 100000)

	n, err := r.Body.Read(b)
	if err == io.EOF {
		log.Printf("EOF")
	} else if err != nil {
		log.Printf("error reading request body: %v", err)
		return
	}

	f, err := ioutil.TempFile("", "edit-server-")
	if err != nil {
		log.Printf("error creating temporary file: %v", err)
		return
	}
	defer f.Close()

	f.Write(b[:n])

	err = exec.Command(*command, f.Name(), url).Run()
	if err != nil {
		log.Printf("error executing edit command %q: %v", *command, err)
		return
	}

	if err := f.Sync(); err != nil {
		log.Printf("error reading temporary file %q: %v", f.Name(), err)
		return
	}

	if _, err := f.Seek(0, 0); err != nil {
		log.Printf("error reading temporary file %q: %v", f.Name(), err)
		return
	}

	n, err = f.Read(b)
	if err != nil {
		log.Printf("error reading temporary file %q: %v", f.Name(), err)
		return
	}

	w.Write(b[:n])
}
