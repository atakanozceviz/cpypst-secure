package controller

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/atakanozceviz/cpypst-secure/model"
)

func PingHandler(w http.ResponseWriter, _ *http.Request) {
	if !Settings.Hidden {
		resp, err := Sign(model.Data{Action: "ping", Content: lname, From: lname, Time: time.Now().Format(time.UnixDate)})
		if err != nil {
			log.Println(err)
			io.WriteString(w, err.Error())
			return
		}
		io.WriteString(w, resp)
	}
	return
}
