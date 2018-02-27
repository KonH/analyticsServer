package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

var history []url.Values

func sendHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	history = append(history, query)
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	response, err := json.Marshal(history)
	if err == nil {
		fmt.Fprintf(w, "%s", response)
	} else {
		log.Fatal(err)
	}
}

func main() {
	fmt.Println("Start server.")
	http.HandleFunc("/send", sendHandler)
	http.HandleFunc("/get", getHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
