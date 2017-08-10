package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/atakanozceviz/cpypst-secure/controller"
	"github.com/atakanozceviz/cpypst-secure/view"
)

const (
	port = "8080"
)

func main() {
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

	http.Handle("/static/", http.FileServer(view.AssetFS()))

	http.HandleFunc("/action", controller.Action)

	http.HandleFunc("/connect", controller.Connect)
	http.HandleFunc("/paste", controller.Paste)
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		controller.Upload(w, r, controller.MPB)
	})

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
