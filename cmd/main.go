package main

import (
	"github.com/KonH/analyticsServer/internal/server"
)

func main() {
	cfg := server.Config{
		ListenTo: ":8080",
		DbHost:   "localhost",
		DbName:   "db",
		CollName: "analytics",
	}
	s := server.New(cfg)
	defer s.Close()
	err := s.Start()
	if err != nil {
		panic(err)
	}
	// TODO: Use wait
}
