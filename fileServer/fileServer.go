package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

var (
	ip   string
	port string
	root string
)

func initialize() {
	flag.StringVar(&ip, "ip", "127.0.0.1", "ip address")
	flag.StringVar(&port, "port", "8080", "port to listen")
	flag.StringVar(&root, "root", "./", "root directory")
}

func main() {
	initialize()
	flag.Parse()
	fmt.Fprintf(os.Stdout,
		"Listening at %s:%s\troot: %s\n", ip, port, root)
	fileServer := http.FileServer(http.Dir(root))
	err := http.ListenAndServe(ip+":"+port, fileServer)

	if err != nil {
		log.Fatal()
		os.Exit(1)
	}
}
