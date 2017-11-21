package controller

import (
	"net/http"
	"regexp"

	"github.com/arschles/go-bindata-html-template"
	"github.com/atakanozceviz/cpypst-secure/model"
	"github.com/atakanozceviz/cpypst-secure/view"
)

var (
	History  model.History
	Incoming model.Connections
	Outgoing model.Connections
	Settings = model.Settings{
		IncomingClip: true,
		IncomingFile: true,
		OutgoingClip: true,
		OutgoingFile: true,
		Hidden:       false,
	}
	checkip = regexp.MustCompile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)
)

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
			Hidden       bool
			Offline      bool
		}{
			Active:       "settings",
			IncomingClip: Settings.IncomingClip,
			IncomingFile: Settings.IncomingFile,
			OutgoingClip: Settings.OutgoingClip,
			OutgoingFile: Settings.OutgoingFile,
			Hidden:       Settings.Hidden,
		}
		tpl.ExecuteTemplate(w, "settings", pData)
	} else {
		w.Write([]byte("You are not authorized!"))
	}
}

func ScanUI(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.New("settings", view.Asset).ParseFiles("view/scan.html", "view/_menu.html"))
	if r.Method == "GET" {
		tpl.ExecuteTemplate(w, "scan", nil)
		return
	}
	servers := model.Connections{}
	// get ip from form
	addr := r.FormValue("ip")
	// check ip if it's valid
	if addr != "" {
		if checkip.Match([]byte(addr)) {
			servers = scan(addr)
		}
	} else {
		servers = scan("")
	}

	pData := struct {
		Active  string
		Servers map[string]*model.Connection
	}{
		Active:  "scan",
		Servers: servers.Connections,
	}
	tpl.ExecuteTemplate(w, "scan", pData)
}
