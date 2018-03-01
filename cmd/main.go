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
	server.Start(cfg)
}
