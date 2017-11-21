package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/atakanozceviz/cpypst-secure/view"
)

func StartApp(port string) error {
	// Get secret key
	for {
		fmt.Print("Enter secret key: ")
		n, _ := fmt.Scanln(&SecretKey)
		if n <= 0 {
			fmt.Println("Secret key cannot be empty")
			continue
		}
		if len(SecretKey) < 3 {
			fmt.Println("Secret key must be longer")
			continue
		}
		break
	}

	// Try to connect automatically
	go func() {
		servers := scan("")
		for _, v := range servers.Connections {
			if err := ConnectTo(v.Ip); err != nil {
				log.Println(err)
			}
		}
	}()

	// Start adding connections
	go func() {
		var addr string
		for {
			for {
				fmt.Println("Enter ip address to add a connection: ")
				n, _ := fmt.Scanln(&addr)
				if n <= 0 {
					fmt.Println("Address cannot be empty")
					continue
				}
				break
			}
			if err := ConnectTo(addr); err != nil {
				log.Println(err)
			}
		}
	}()

	http.HandleFunc("/", HistoryUI)
	http.HandleFunc("/connections", ConnectionsUI)
	http.HandleFunc("/settings", SettingsUI)
	http.HandleFunc("/scan", ScanUI)
	http.Handle("/static/", http.FileServer(view.AssetFS()))

	http.HandleFunc("/action", ActionHandler)

	http.HandleFunc("/ping", PingHandler)
	http.HandleFunc("/connect", ConnectHandler)
	http.HandleFunc("/paste", PasteHandler)
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		UploadHandler(w, r, MPB)
	})
	return http.ListenAndServe(":"+port, nil)
}
