package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/globalsign/mgo"
)

// Public structs

// Config for server settings
type Config struct {
	ListenTo string
	DbHost   string
	DbName   string
	CollName string
}

// Server to handle actions
type Server struct {
	cfg Config
	srv *http.Server
	db  *mgo.Session
}

// Public funcs

// New server creation with given config
func New(config Config) Server {
	return Server{
		cfg: config,
	}
}

// Start server
func (s *Server) Start() error {
	s.addHandlers()
	err := s.connectToDb()
	if err != nil {
		return err
	}
	s.startServer()
	return nil
}

// Stop server
func (s *Server) Stop() {
	if s.srv != nil {
		s.srv.Shutdown(nil)
	}
}

// Close server
func (s *Server) Close() {
	if s.db != nil {
		s.db.Close()
	}
}

// Internals

func (s *Server) addHandlers() {
	http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		s.handleSend(w, r)
	})
	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		s.handleGet(w, r)
	})
}

func (s *Server) connectToDb() error {
	fmt.Println("Connect to mongodb.")
	db, err := mgo.Dial(s.cfg.DbHost)
	s.db = db
	return err
}

func (s *Server) startServer() {
	fmt.Println("Start server.")
	s.srv = &http.Server{Addr: s.cfg.ListenTo}
	go func() {
		if err := s.srv.ListenAndServe(); err != nil {
			log.Printf("HttpServer: ListenAndServe() error: %s", err)
		}
	}()
}

func (s *Server) handleSend(w http.ResponseWriter, r *http.Request) {
	item := r.URL.Query()
	err := s.addItem(item)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Server) addItem(item map[string][]string) error {
	c := s.db.DB(s.cfg.DbName).C(s.cfg.CollName)
	return c.Insert(item)
}

func (s *Server) handleGet(w http.ResponseWriter, r *http.Request) {
	items, err := s.getItems()
	if err != nil {
		log.Fatal(err)
	}
	jsonResult, err := json.Marshal(items)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "%s", jsonResult)
}

func (s *Server) getItems() ([]interface{}, error) {
	c := s.db.DB(s.cfg.DbName).C(s.cfg.CollName)
	var result []interface{}
	err := c.Find(nil).All(&result)
	return result, err
}
