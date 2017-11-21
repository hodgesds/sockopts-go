package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/hodgesds/sockopts-go"
)

var pid = os.Getpid()

type helloHandler struct{}

func (h helloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("serving from %d\n", pid)
	fmt.Fprintf(w, "hello, from pid %d\n", pid)
}

func main() {
	listener, err := sockopts.SockoptListener(
		"tcp",
		":8080",
		sockopts.SO_REUSEPORT,
		sockopts.TCP_FASTOPEN,
	)
	if err != nil {
		log.Fatalln(err)
	}
	http.Serve(listener, helloHandler{})
}
