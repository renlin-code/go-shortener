package main

import (
	"github.com/renlin-code/go-shortener/internal/config"
)

func main() {
	cfg := config.MustLoad()

	//TODO: init logger (slog)

	//TODO: init storage (SQL Lite)

	//TODO: init router (chi)

	//TODO: run server
}
