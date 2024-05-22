package main

import (
	"github.com/mathesukkj/gourl-shortener/api"
	"github.com/mathesukkj/gourl-shortener/sqlite"
)

func main() {
	server := api.NewServer()
	db := sqlite.NewDB("./sqlite/database/url_shortener")
	server.URLService = sqlite.NewURLService(db)
	server.Serve(":4090")
}
