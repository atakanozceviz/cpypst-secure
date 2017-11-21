package main

import (
	"flag"
	"log"

	"github.com/atakanozceviz/cpypst-secure/controller"
)

var port = flag.String("port", "8080", "Specify the port you want to use")

func main() {
	controller.StartApp(*port)
}

func init() {
	// Enable line numbers in logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	flag.Parse()
	controller.Port = *port
}
