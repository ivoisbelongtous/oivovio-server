package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// Don't save API key in source.
var apiKey string

type ScrobbleStatus int

const (
	Ready ScrobbleStatus = iota
	Fetching
)

type ScrobbleRequest struct {
	status ScrobbleStatus
	json   string
}

// Request REST response from last.fm API.
func GetScrobble() string {
	res, err := http.Get("http://ws.audioscrobbler.com/2.0/?" +
		"method=user.getRecentTracks&" +
		"api_key=" + apiKey +
		"&user=Doomboy95" +
		"&limit=1" +
		"&format=json")
	if err != nil {
		log.Fatal(err)
	}

	var buf bytes.Buffer
	io.Copy(&buf, res.Body)
	res.Body.Close()
	return buf.String()
}

func ScrobbleMonitor(out chan ScrobbleRequest) {
	// Cache times out after 30 seconds
	timeout := time.After(30 * time.Second)

	sr := ScrobbleRequest{Fetching, ""}
	for {
		select {
		case out <- sr:
			if sr.status == Fetching {
				sr = ScrobbleRequest{Ready, GetScrobble()}
				// Reset timeout
				timeout = time.After(30 * time.Second)
			}
		case <-timeout:
			// Invalidate the cache
			sr = ScrobbleRequest{Fetching, ""}
		}
	}
}

func main() {
	r, err := os.Open("apikey.txt")
	if err != nil {
		log.Fatal(err)
	}
	_, err = fmt.Fscan(r, &apiKey)
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan ScrobbleRequest)
	go ScrobbleMonitor(c)

	http.HandleFunc("/scrobbles.json",
		func(w http.ResponseWriter, req *http.Request) {
			// Poll until valid scrobble returned
			for {
				r := <-c
				if r.status == Ready {
					io.WriteString(w, r.json)
					return
				}
			}
		})
	certRoot := "/etc/letsencrypt/live/oivov.io/"
	http.ListenAndServeTLS(
		":8345",
		certRoot+"cert.pem",
		certRoot+"privkey.pem",
		nil)
}
