package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	jsonStr, err := json.Marshal(query)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "%s", jsonStr)
}

func main() {
	fmt.Println("Start server.")
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
