package controller

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/atakanozceviz/cpypst-secure/model"
)

func scan(addr string) model.Connections {
	servers := model.Connections{}
	var ips []string

	if addr == "" {
		auto, err := localIP()
		if err != nil {
			log.Println(err)
			return servers
		}
		ips = auto
	} else {
		// delete last n characters
		addr = strings.TrimRightFunc(addr, isNumber)
		ips = append(ips, addr)
	}

	if len(ips) == 1 {
		var wg sync.WaitGroup
		wg.Add(255)

		client := http.Client{
			Timeout: time.Duration(time.Second * 2),
		}
		// scan the network
		for i := 1; i <= 255; i++ {
			ip := ips[0] + strconv.Itoa(i)
			go func(ip string) {
				req, _ := http.NewRequest(http.MethodGet, "http://"+ip+":"+Port+"/ping", nil)
				resp, err := client.Do(req)
				if err == nil {
					// read the response and add to "servers"
					bdy, _ := ioutil.ReadAll(resp.Body)
					defer resp.Body.Close()
					data, err := Parse(string(bdy))
					if err == nil && data["Action"] == "ping" && data["From"] != lname {
						servers.Add(model.Connection{ip, data["From"].(string), true, time.Now().Format(time.UnixDate)})
					}
				}
				wg.Done()
			}(ip)
		}
		wg.Wait()
	}
	return servers
}

func localIP() ([]string, error) {
	var ips []string

	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		addrs, err := net.InterfaceAddrs()
		if err != nil {
			return nil, err
		}
		for _, a := range addrs {
			if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					ip := ipnet.IP.String()
					ip = strings.TrimRightFunc(ip, isNumber)
					ips = append(ips, ip)
				}
			}
		}
		return ips, nil
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr).IP.String()
	ips = append(ips, strings.TrimRightFunc(localAddr, isNumber))
	return ips, nil
}

func isNumber(r rune) bool {
	return unicode.IsNumber(r)
}
