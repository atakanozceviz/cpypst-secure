package controller

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/atakanozceviz/cpypst-secure/model"
	"github.com/atotto/clipboard"
	"github.com/vbauerster/mpb/decor"
	"gopkg.in/vbauerster/mpb.v3"
)

var tmp model.Tmp
var clip model.Tmp
var lname, _ = os.Hostname()

var re = regexp.MustCompile(`:[0-9]+`)

func Connect(w http.ResponseWriter, r *http.Request) {
	addr := r.RemoteAddr
	defer r.Body.Close()
	rbody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	data, err := Parse(string(rbody))
	if err != nil {
		log.Println(err)
		io.WriteString(w, "Wrong request!")
		return
	}
	if data["Action"] == "connect" {
		name := data["From"].(string)
		ip := re.ReplaceAllString(addr, "")
		Incoming.Add(model.Connection{Ip: ip, Name: name, Active: true, Time: time.Now().Format(time.UnixDate)})
		w.Write([]byte(lname))
		fmt.Println("\n" + name + " (" + ip + ") is connected!")
	}
}

func Paste(_ http.ResponseWriter, r *http.Request) {
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

func Upload(w http.ResponseWriter, r *http.Request, p *mpb.Progress) {
	if Settings.IncomingFile == true {
		addr := r.RemoteAddr
		ip := re.ReplaceAllString(addr, "")

		// get file size and name of the data
		filename := r.Header.Get("FileName")
		filesize, err := strconv.ParseInt(r.Header.Get("FileSize"), 10, 64)
		if err != nil {
			log.Println(err)
		}

		if r.Method == "POST" && filename != "" && err == nil {
			// create tmp folder if not exists
			if _, err := os.Stat("./tmp"); os.IsNotExist(err) {
				os.Mkdir("tmp", 0777)
			}
			// create file to write received data
			dest, err := os.Create("./tmp/" + filename)
			if err != nil {
				log.Println(err)
			}
			defer dest.Close()

			// create and start bar
			bar := p.AddBar(filesize,
				mpb.PrependDecorators(
					decor.StaticName(filename+"(receive)", 0, decor.DwidthSync|decor.DidentRight),
					decor.Counters("%3s / %3s", decor.Unit_kB, 18, decor.DSyncSpace),
				),
				mpb.AppendDecorators(decor.Percentage(5, 0)),
			)

			// create proxy reader
			reader := bar.ProxyReader(r.Body)
			defer r.Body.Close()

			// copy received data to destination
			n, err := io.Copy(dest, reader)
			if err != nil {
				log.Println(err)
			}
			w.Write([]byte(fmt.Sprintf("%d bytes are recieved.\n", n)))

			// remove bar
			p.RemoveBar(bar)

			dir, err := filepath.Abs("./tmp/" + filename)
			if err != nil {
				log.Println(err)
			}

			clip.Write(dir)
			tmp.Write(dir)

			// add to history
			History.Add(model.HistItem{Content: dir, Time: time.Now().Format(time.UnixDate), Ip: ip})
		}
	} else {
		w.Write([]byte(lname + " disabled incoming file!"))
	}
}

func Action(w http.ResponseWriter, r *http.Request) {
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
		default:
			io.WriteString(w, "Invalid action!")
		}
	} else {
		io.WriteString(w, "You are not authorized!")
	}
}
