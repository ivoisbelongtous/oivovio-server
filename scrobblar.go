package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var apiKey string
var scrobbleJSON string

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
	scrobbleJSON = buf.String()
	return scrobbleJSON
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
