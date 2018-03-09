package server

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			s.handleSend(w, r)
		} else {
			s.handleGet(w, r)
		}
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
	item, err := parseBody(r.Body)
	if err != nil {
		log.Println(err)
		return
	}
	err = s.addItem(item)
	if err != nil {
		log.Println(err)
	}
}

func parseBody(bodyReader io.ReadCloser) (interface{}, error) {
	body, err := ioutil.ReadAll(bodyReader)
	if err != nil {
		return nil, err
	}
	log.Printf("parseBody: body = '%s'", body)
	var item interface{}
	err = json.Unmarshal(body, &item)
	if err != nil {
		return nil, err
	}
	log.Printf("parseBody: item = '%s'", item)
	return item, nil
}

func (s *Server) addItem(item interface{}) error {
	c := s.db.DB(s.cfg.DbName).C(s.cfg.CollName)
	return c.Insert(item)
}

func (s *Server) handleGet(w http.ResponseWriter, r *http.Request) {
	items, err := s.getItems()
	if err != nil {
		log.Println(err)
		return
	}
	jsonResult, err := json.Marshal(items)
	if err != nil {
		log.Println(err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", jsonResult)
}

func (s *Server) getItems() ([]interface{}, error) {
	c := s.db.DB(s.cfg.DbName).C(s.cfg.CollName)
	var result []interface{}
	err := c.Find(nil).All(&result)
	return result, err
}

func (s *Server) clearItems() error {
	c := s.db.DB(s.cfg.DbName).C(s.cfg.CollName)
	return c.DropCollection()
}
