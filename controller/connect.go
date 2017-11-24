package controller

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/atakanozceviz/cpypst-secure/model"
)

func ConnectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		addr := r.RemoteAddr
		defer r.Body.Close()
		rbody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			return
		}

		data, err := Parse(string(rbody))
		if err != nil {
			log.Println(err)
			io.WriteString(w, "Cannot parse!")
			return
		}

		if data["Action"] == "connect" {
			name := data["From"].(string)
			ip := re.ReplaceAllString(addr, "")
			Incoming.Add(model.Connection{Ip: ip, Name: name, Active: true, Time: time.Now().Format(time.UnixDate)})

			resp, err := Sign(model.Data{Action: "name", Content: lname, From: lname, Time: time.Now().Format(time.UnixDate)})
			if err != nil {
				log.Println(err)
			}
			w.Write([]byte(resp))
			// Connect to who is connected
			if !Outgoing.Has(ip) {
				Outgoing.Add(model.Connection{Ip: ip, Name: name, Active: true, Time: time.Now().Format(time.UnixDate)})
				if err := ConnectTo(ip); err != nil {
					log.Println(err)
				}
			}

			fmt.Printf("\n%s (%s) is connected!\n", name, ip)
			return
		}
	}
	io.WriteString(w, "Wrong request!")
}
