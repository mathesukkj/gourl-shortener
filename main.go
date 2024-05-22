package main

import "github.com/mathesukkj/gourl-shortener/api"

func main() {
	server := api.NewServer()
	server.Serve(":4090")
}
