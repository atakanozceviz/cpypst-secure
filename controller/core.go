package controller

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/atakanozceviz/cpypst-secure/model"
	"github.com/atotto/clipboard"
	"github.com/vbauerster/mpb/decor"
	"gopkg.in/vbauerster/mpb.v3"
)

var MPB = mpb.New(
	mpb.WithWidth(64),
	mpb.WithFormat("╢▌▌░╟"),
)

func checkClip() {
	var x string
	var err error
	for {
		x, err = clipboard.ReadAll()
		if err != nil && err.Error() != "exit status 1" {
			log.Println(err)
		}
		if x != "" {
			clip.Write(x)
			if clip.Read() != tmp.Read() {
				tmp.Write(clip.Read())
				send(clip.Read())
			}
		}
		time.Sleep(time.Second * 1)
	}
}

func send(clip string) {
	addr := Outgoing.Connections

	if len(addr) > 0 {
		for _, v := range addr {
			if v.Active == true {
				if ok := path.IsAbs(clip); ok && Settings.OutgoingFile == true {
					go postFile(clip, "http://"+v.Ip+":8080/upload", MPB)
				} else if Settings.OutgoingClip == true {
					_, err := EncSend("paste", lname, clip, v.Ip)
					if err != nil {
						log.Println(err)
					}
				}
			}
		}
	}
}

func postFile(fp, url string, p *mpb.Progress) {
	fp, _ = filepath.Abs(fp)
	re, err := regexp.Compile("[\r\n]")
	if err != nil {
		log.Println(err)
	}

	fp = re.ReplaceAllString(fp, "")

	file, err := os.Open(fp)
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		log.Println(err)
		return
	}
	if info.IsDir() {
		fmt.Println("Cannot send directory")
		return
	}
	// create bar
	bar := p.AddBar(info.Size(),
		mpb.PrependDecorators(
			decor.StaticName(info.Name()+"(send)", len(info.Name()), decor.DwidthSync|decor.DidentRight),
			decor.Counters("%3s / %3s", decor.Unit_kB, 18, decor.DSyncSpace),
		),
		mpb.AppendDecorators(decor.Percentage(5, 0)),
	)

	// create multi writer
	writer := bar.ProxyReader(file)

	req, err := http.NewRequest(http.MethodPost, url, writer)
	if err != nil {
		log.Println(err)
		return
	}
	req.Header.Set("Content-Type", "binary/octet-stream")
	req.Header.Set("FileName", info.Name())
	req.Header.Set("FileSize", strconv.FormatInt(info.Size(), 10))

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()
	message, _ := ioutil.ReadAll(res.Body)

	if len(message) > 0 {
		log.Println(string(message))
	}
	// remove bar
	p.RemoveBar(bar)

}

func Start() {
	var addr string
	for {
		for {
			fmt.Print("Enter ip address to add a connection: ")
			n, _ := fmt.Scanln(&addr)
			if n <= 0 {
				fmt.Println("Address cannot be empty")
				continue
			}
			break
		}

		resp, err := EncSend("connect", lname, "", addr)
		if err != nil {
			log.Println(err)
			continue
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			resp.Body.Close()
			continue
		}

		Outgoing.Add(model.Connection{Ip: addr, Name: string(body), Active: true, Time: time.Now().Format(time.UnixDate)})
		resp.Body.Close()
		if len(Outgoing.Connections) == 1 {
			go checkClip()
		}
	}
}
