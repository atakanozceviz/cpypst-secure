package controller

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/atakanozceviz/cpypst-secure/model"
	"github.com/atotto/clipboard"
)

func PasteHandler(_ http.ResponseWriter, r *http.Request) {
	ip := re.ReplaceAllString(r.RemoteAddr, "")
	if Incoming.Connections[ip].Active == true && Settings.IncomingClip == true {
		defer r.Body.Close()
		rbody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
		}
		data, err := Parse(string(rbody))
		if err != nil {
			log.Println(err)
		} else if data["Action"] == "paste" {
			clip.Write(data["Content"].(string))
			tmp.Write(data["Content"].(string))
			if clipboard.WriteAll(clip.Read()) != nil {
				log.Println(err)
			}
			History.Add(model.HistItem{Ip: ip, Content: clip.Read(), Time: data["Time"].(string)})
		}
	}
}
