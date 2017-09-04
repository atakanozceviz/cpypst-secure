package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/atakanozceviz/cpypst-secure/controller"
	"github.com/atakanozceviz/cpypst-secure/view"
)

var port = flag.String("port", "8080", "Specify the port you want to use")

func init() {
	controller.Port = *port
}

func main() {
	flag.Parse()

	for {
		fmt.Print("Enter secret key: ")
		n, _ := fmt.Scanln(&controller.SecretKey)
		if n <= 0 {
			fmt.Println("Secret key cannot be empty")
			continue
		}
		if len(controller.SecretKey) < 3 {
			fmt.Println("Secret key must be longer")
			continue
		}
		break
	}

	go controller.Start()

	http.HandleFunc("/", controller.HistoryUI)
	http.HandleFunc("/connections", controller.ConnectionsUI)
	http.HandleFunc("/settings", controller.SettingsUI)
	http.HandleFunc("/scan", controller.ScanUI)
	http.Handle("/static/", http.FileServer(view.AssetFS()))

	http.HandleFunc("/action", controller.Action)

	http.HandleFunc("/ping", controller.Ping)
	http.HandleFunc("/connect", controller.Connect)
	http.HandleFunc("/paste", controller.Paste)
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		controller.Upload(w, r, controller.MPB)
	})

	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

func init() {
	// Enable line numbers in logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
