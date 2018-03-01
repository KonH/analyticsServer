package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/globalsign/mgo"
)

type dbConfig struct {
	dbName   string
	collName string
}

func sendHandler(cfg dbConfig, db *mgo.Session, w http.ResponseWriter, r *http.Request) {
	c := db.DB(cfg.dbName).C(cfg.collName)
	item := r.URL.Query()
	c.Insert(item)
}

func getHandler(cfg dbConfig, db *mgo.Session, w http.ResponseWriter, r *http.Request) {
	c := db.DB(cfg.dbName).C(cfg.collName)
	var result []interface{}
	err := c.Find(nil).All(&result)
	if err != nil {
		log.Fatal(err)
	}
	jsonResult, err := json.Marshal(result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "%s", jsonResult)
}

func main() {
	cfg := dbConfig{"db", "analytics"}

	fmt.Println("Connect to mongodb.")
	db, err := mgo.Dial("localhost")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("Start server.")
	http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		sendHandler(cfg, db, w, r)
	})
	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		getHandler(cfg, db, w, r)
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
