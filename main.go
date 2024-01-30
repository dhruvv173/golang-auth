package main

import (
	"log"

	"github.com/dhruvv173/auth/db"
	"github.com/dhruvv173/auth/handlers"
	"github.com/dhruvv173/auth/session"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func main() {
	if err:= db.InitDB(); err != nil{
		log.Fatal("Error initializing DB:", err)
	}
	defer db.CloseDB()
	session.Init()
	
	router := gin.Default()
	router.Use(sessions.Sessions("qid", session.GetStore()))
	router.GET("/ok", handlers.OKhandler)
	router.POST("/register", handlers.RegisterHandler)
	router.POST("/login", handlers.LoginHandler)
	router.POST("/logout", handlers.LogoutHandler)
	
	router.GET("/protected", handlers.IsAuth(), handlers.CurrentUserHandler)

	if err:= router.Run(); err != nil{
		log.Fatal("Error starting server:", err)
	}
}