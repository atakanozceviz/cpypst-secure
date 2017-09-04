package controller

import (
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/arschles/go-bindata-html-template"
	"github.com/atakanozceviz/cpypst-secure/model"
	"github.com/atakanozceviz/cpypst-secure/view"
)

var checkip = regexp.MustCompile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)

var History model.History
var Incoming model.Connections
var Outgoing model.Connections
var Settings = model.Settings{
	IncomingClip: true,
	IncomingFile: true,
	OutgoingClip: true,
	OutgoingFile: true,
	Hidden:       false,
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
	servers := model.Connections{}
	// get ip from form
	fip := r.FormValue("ip")
	// check ip if it's valid
	if checkip.Match([]byte(fip)) {
		// delete last n characters
		fip = fip[:len(fip)-len(checkip.FindStringSubmatch(fip)[3])]

		var wg sync.WaitGroup
		wg.Add(255)

		client := http.Client{
			Timeout: time.Duration(time.Second * 2),
		}
		// scan the network
		for i := 1; i <= 255; i++ {
			ip := fip + strconv.Itoa(i)
			go func(ip string) {
				req, _ := http.NewRequest(http.MethodGet, "http://"+ip+":"+Port+"/ping", nil)
				resp, err := client.Do(req)
				if err == nil {
					// read the response and add to "servers"
					name, err := ioutil.ReadAll(resp.Body)
					if err == nil && len(name) > 0 {
						servers.Add(model.Connection{ip, string(name), true, time.Now().Format(time.UnixDate)})
					}
					resp.Body.Close()
				}
				wg.Done()
			}(ip)
		}
		wg.Wait()
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
