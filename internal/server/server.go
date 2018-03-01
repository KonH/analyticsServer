package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/globalsign/mgo"
)

// Config
type Config struct {
	ListenTo string
	DbHost   string
	DbName   string
	CollName string
}

func sendHandler(cfg *Config, db *mgo.Session, w http.ResponseWriter, r *http.Request) {
	c := db.DB(cfg.DbName).C(cfg.CollName)
	item := r.URL.Query()
	c.Insert(item)
}

func getHandler(cfg *Config, db *mgo.Session, w http.ResponseWriter, r *http.Request) {
	c := db.DB(cfg.DbName).C(cfg.CollName)
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

// Start server with given config
func Start(cfg Config) {
	fmt.Println("Connect to mongodb.")
	db, err := mgo.Dial(cfg.DbHost)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("Start server.")
	http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		sendHandler(&cfg, db, w, r)
	})
	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		getHandler(&cfg, db, w, r)
	})
	log.Fatal(http.ListenAndServe(cfg.ListenTo, nil))
}
