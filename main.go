package main

import (
	_ "github.com/joho/godotenv/autoload"

	"github.com/mathesukkj/gourl-shortener/api"
	"github.com/mathesukkj/gourl-shortener/redis"
	"github.com/mathesukkj/gourl-shortener/sqlite"
)

func main() {
	server := api.NewServer()
	db := sqlite.NewDB("./sqlite/database/url_shortener")
	server.URLService = sqlite.NewURLService(db)
	server.RedisClient = redis.New()
	server.Serve(":4090")
}
