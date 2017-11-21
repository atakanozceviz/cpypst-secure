package controller

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/atakanozceviz/cpypst-secure/model"
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
)

func UploadHandler(w http.ResponseWriter, r *http.Request, p *mpb.Progress) {
	if Settings.IncomingFile == true {
		addr := r.RemoteAddr
		ip := re.ReplaceAllString(addr, "")

		// Get file size and name of the data
		filename := r.Header.Get("FileName")
		filesize, err := strconv.ParseInt(r.Header.Get("FileSize"), 10, 64)
		if err != nil {
			log.Println(err)
		}

		if r.Method == "POST" && filename != "" && err == nil {
			// Create tmp folder if not exists
			if _, err := os.Stat("./tmp"); os.IsNotExist(err) {
				os.Mkdir("tmp", 0777)
			}
			// Create file to write received data
			dest, err := os.Create("./tmp/" + filename)
			if err != nil {
				log.Println(err)
			}
			defer dest.Close()

			// Create and start bar
			bar := p.AddBar(filesize,
				mpb.PrependDecorators(
					decor.StaticName(filename+"(receive)", 0, decor.DwidthSync|decor.DidentRight),
					decor.Counters("%3s / %3s", decor.Unit_kB, 18, decor.DSyncSpace),
				),
				mpb.AppendDecorators(decor.Percentage(5, 0)),
			)
			// Remove bar
			defer p.RemoveBar(bar)

			// Create proxy reader
			reader := bar.ProxyReader(r.Body)
			defer r.Body.Close()

			// Copy received data to destination
			n, err := io.Copy(dest, reader)
			if err != nil {
				log.Println(err)
			}
			w.Write([]byte(fmt.Sprintf("%d bytes are recieved.\n", n)))

			dir, err := filepath.Abs("./tmp/" + filename)
			if err != nil {
				log.Println(err)
			}

			clip.Write(dir)
			tmp.Write(dir)

			// Add to history
			History.Add(model.HistItem{Content: dir, Time: time.Now().Format(time.UnixDate), Ip: ip})
		}
	} else {
		w.Write([]byte(lname + " disabled incoming file!"))
	}
}
