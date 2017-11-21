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
			w.Write([]byte(lname))
			fmt.Println("\n" + name + " (" + ip + ") is connected!")
			return
		}
	}
	io.WriteString(w, "Wrong request!")
}
