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
	"strings"
	"time"

	"github.com/atakanozceviz/cpypst-secure/model"
	"github.com/atotto/clipboard"
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
)

var (
	MPB = mpb.New(
		mpb.WithWidth(64),
		mpb.WithFormat("╢▌▌░╟"),
	)

	tmp      model.Tmp
	clip     model.Tmp
	lname, _ = os.Hostname()
	re       = regexp.MustCompile(`:[0-9]+`)
	winPath  = regexp.MustCompile(`^"?([a-zA-Z]:|\\\\[^/\\:*?"<>|]+\\[^/\\:*?"<>|]+)(\\[^/\\:*?"<>|]+)+(\.[^/\\:*?"<>|]+)"?$`)
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
				if ok := path.IsAbs(clip); ok || winPath.Match([]byte(clip)) && Settings.OutgoingFile == true {
					clip = strings.TrimSpace(clip)
					clip = strings.Trim(clip, "\"")
					go postFile(clip, "http://"+v.Ip+":"+Port+"/upload", MPB)
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
	// Create bar
	fn := ""
	if len(info.Name()) > 25 {
		fn = info.Name()[:20] + "..."
	} else {
		fn = info.Name()
	}
	bar := p.AddBar(info.Size(),
		mpb.PrependDecorators(
			decor.StaticName(fn+"(send)", len(info.Name()), decor.DwidthSync|decor.DidentRight),
			decor.Counters("%3s / %3s", decor.Unit_kB, 18, decor.DSyncSpace),
		),
		mpb.AppendDecorators(decor.Percentage(5, 0)),
	)
	// Remove bar
	defer p.RemoveBar(bar)

	// Create multi writer
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
}

func ConnectTo(addr string) error {
	resp, err := EncSend("connect", lname, "", addr)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	Outgoing.Add(model.Connection{Ip: addr, Name: string(body), Active: true, Time: time.Now().Format(time.UnixDate)})

	if len(Outgoing.Connections) == 1 {
		go checkClip()
	}
	return nil
}
