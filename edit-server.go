package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
)

var (
	port    = flag.Int("port", 8989, "port to run the server on")
	command = flag.String("command", "~/edit.sh", "command to run to open the editor. Use '%s' as placeholder for the temporary file.")
)

func main() {
	flag.Parse()

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
		panic(err)
	}

	f, err := ioutil.TempFile("", "edit-server-")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.Write(b[:n])

	err = exec.Command(*command, f.Name()).Run()
	if err != nil {
		panic(err)
	}

	if err := f.Sync(); err != nil {
		panic(err)
	}

	if _, err := f.Seek(0, 0); err != nil {
		panic(err)
	}

	n, err = f.Read(b)
	if err != nil {
		panic(err)
	}

	w.Write(b[:n])
}
