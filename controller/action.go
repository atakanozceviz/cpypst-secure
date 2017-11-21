package controller

import (
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/atotto/clipboard"
)

func ActionHandler(w http.ResponseWriter, r *http.Request) {
	if re.ReplaceAllString(r.RemoteAddr, "") == "127.0.0.1" {
		switch action := r.FormValue("action"); action {
		case "dlt":
			{
				id, err := strconv.Atoi(r.FormValue("ID"))
				if err != nil {
					log.Println(err)
				} else {
					History.Remove(id)
				}
			}
		case "cpy":
			{
				id, err := strconv.Atoi(r.FormValue("ID"))
				if err != nil {
					log.Println(err)
				} else {
					cpy := History.History[id].Content
					tmp.Write(cpy)
					clipboard.WriteAll(cpy)
				}
			}
		case "ienable":
			{
				id := r.FormValue("ID")
				if id != "" {
					if val, ok := Incoming.Connections[id]; ok {
						val.Active = true
					} else {
						io.WriteString(w, "Couldn't find connection")
					}
				}
			}
		case "idisable":
			{
				id := r.FormValue("ID")
				if id != "" {
					if val, ok := Incoming.Connections[id]; ok {
						val.Active = false
					} else {
						io.WriteString(w, "Couldn't find connection")
					}
				}
			}
		case "oenable":
			{
				id := r.FormValue("ID")
				if id != "" {
					if val, ok := Outgoing.Connections[id]; ok {
						val.Active = true
					} else {
						io.WriteString(w, "Couldn't find connection")
					}
				}
			}
		case "odisable":
			{
				id := r.FormValue("ID")
				if id != "" {
					if val, ok := Outgoing.Connections[id]; ok {
						val.Active = false
					} else {
						io.WriteString(w, "Couldn't find connection")
					}
				}
			}
		case "incomingclip":
			{
				if Settings.IncomingClip == true {
					Settings.IncomingClip = false
				} else {
					Settings.IncomingClip = true
				}
			}
		case "outgoingclip":
			{
				if Settings.OutgoingClip == true {
					Settings.OutgoingClip = false
				} else {
					Settings.OutgoingClip = true
				}
			}
		case "incomingfile":
			{
				if Settings.IncomingFile == true {
					Settings.IncomingFile = false
				} else {
					Settings.IncomingFile = true
				}
			}
		case "outgoingfile":
			{
				if Settings.OutgoingFile == true {
					Settings.OutgoingFile = false
				} else {
					Settings.OutgoingFile = true
				}
			}
		case "connect":
			{
				ip := r.FormValue("ID")
				if ip != "" {
					if err := ConnectTo(ip); err != nil {
						log.Println(err)
						return
					}
					io.WriteString(w, "success")
				}
			}
		case "hidden":
			{
				if Settings.Hidden == true {
					Settings.Hidden = false
				} else {
					Settings.Hidden = true
				}
			}
		default:
			io.WriteString(w, "Invalid action!")
		}
	} else {
		io.WriteString(w, "You are not authorized!")
	}
}
