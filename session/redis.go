package session

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/joho/godotenv"
)

var store sessions.Store

func Init(){
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	secret := os.Getenv("SECRET")
	store, err = redis.NewStore(10,"tcp", "127.0.0.1:6379","", []byte(secret))
	if err != nil {
		log.Fatalf("Error creating Redis store: %v", err)
	}
	store.Options(sessions.Options{
		Path: "/",
		MaxAge: 3600,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	
}

func GetStore() sessions.Store{
	return store
}