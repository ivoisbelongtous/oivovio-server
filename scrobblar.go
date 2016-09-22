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

var apiKey string

type ScrobbleRequest struct {
	json        string
	requestTime time.Time
}

var scrobble ScrobbleRequest

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

	scrobble.json = buf.String()
	resTime := res.Header.Get("Date")
	scrobble.requestTime, err = time.Parse(time.RFC1123, resTime)
	if err != nil {
		log.Fatal(err)
	}
	return scrobble.json
}

func ScrobbleServer(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, GetScrobble())
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

	http.HandleFunc("/scrobbles.json", ScrobbleServer)
	http.ListenAndServe(":8345", nil)
}
