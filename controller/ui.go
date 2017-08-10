package controller

import (
	"net/http"

	"github.com/arschles/go-bindata-html-template"
	"github.com/atakanozceviz/cpypst-secure/model"
	"github.com/atakanozceviz/cpypst-secure/view"
)

var History model.History
var Incoming model.Connections
var Outgoing model.Connections
var Settings = model.Settings{
	IncomingClip: true,
	IncomingFile: true,
	OutgoingClip: true,
	OutgoingFile: true,
}

func HistoryUI(w http.ResponseWriter, r *http.Request) {
	if re.ReplaceAllString(r.RemoteAddr, "") == "127.0.0.1" {
		tpl := template.Must(template.New("history", view.Asset).ParseFiles("view/history.html", "view/_menu.html"))
		pData := struct {
			Active   string
			History  []model.HistItem
			Incoming map[string]*model.Connection
		}{
			Active:   "history",
			History:  History.History,
			Incoming: Incoming.Connections,
		}

		tpl.ExecuteTemplate(w, "history", pData)
	} else {
		w.Write([]byte("You are not authorized!"))
	}
}

func ConnectionsUI(w http.ResponseWriter, r *http.Request) {
	if re.ReplaceAllString(r.RemoteAddr, "") == "127.0.0.1" {
		tpl := template.Must(template.New("connections", view.Asset).ParseFiles("view/connections.html", "view/_menu.html"))
		pData := struct {
			Active   string
			Incoming map[string]*model.Connection
			Outgoing map[string]*model.Connection
		}{
			Active:   "connections",
			Incoming: Incoming.Connections,
			Outgoing: Outgoing.Connections,
		}

		tpl.ExecuteTemplate(w, "connections", pData)
	} else {
		w.Write([]byte("You are not authorized!"))
	}
}

func SettingsUI(w http.ResponseWriter, r *http.Request) {
	if re.ReplaceAllString(r.RemoteAddr, "") == "127.0.0.1" {
		tpl := template.Must(template.New("settings", view.Asset).ParseFiles("view/settings.html", "view/_menu.html"))
		pData := struct {
			Active       string
			IncomingClip bool
			IncomingFile bool
			OutgoingClip bool
			OutgoingFile bool
		}{
			Active:       "settings",
			IncomingClip: Settings.IncomingClip,
			IncomingFile: Settings.IncomingFile,
			OutgoingClip: Settings.OutgoingClip,
			OutgoingFile: Settings.OutgoingFile,
		}
		tpl.ExecuteTemplate(w, "settings", pData)
	} else {
		w.Write([]byte("You are not authorized!"))
	}
}
