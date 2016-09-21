package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var apiKey string

func ScrobbleServer(w http.ResponseWriter, req *http.Request) {
	res, err := http.Get("http://ws.audioscrobbler.com/2.0/?" +
		"method=user.getRecentTracks&" +
		"api_key=" + apiKey +
		"&user=Doomboy95" +
		"&limit=1" +
		"&format=json")
	if err != nil {
		log.Fatal(err)
	}
	io.Copy(w, res.Body)
	res.Body.Close()
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

	http.HandleFunc("/scrobbles", ScrobbleServer)
	http.ListenAndServe(":8345", nil)
}
