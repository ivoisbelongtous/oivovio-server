package main

import (
	"io"
	"net/http"
)

func ScrobbleServer(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Listen to more music!\n")
}

func main() {
	http.HandleFunc("/scrobbles", ScrobbleServer)
	http.ListenAndServe(":8345", nil)
}
