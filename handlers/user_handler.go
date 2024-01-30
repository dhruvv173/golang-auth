package handlers

import (
	"log"
	"net/http"

	"github.com/dhruvv173/auth/db"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(c *gin.Context){

	name := c.PostForm("name")
	email := c.PostForm("email")
	password := c.PostForm("password")
	confirmedPassword := c.PostForm("confirmedPassword")

	if name == "" || email == "" || password == "" || confirmedPassword == ""{
		c.JSON(400, gin.H{"error":"All fields are required"})
		return
	}
	if password != confirmedPassword {
		c.JSON(400, gin.H{"error":"Passwords dont  match"})
		return
	}
	hashedPassword, err := hashPassword(password)
	if err != nil {
		c.JSON(500, gin.H{"error":"Internal server error"})
		return
	}
	newUser := db.User{
		Name: name,
		Email: email,
		Password: hashedPassword,
	}
	dbInstance := db.GetDB()
	if dbInstance == nil {
		log.Fatal("GORM database instance is nil")
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
        return
	}
	result := dbInstance.Create(&newUser)
	if result.Error != nil {
		log.Fatal("Error signing up:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
        return
	}
	c.JSON(201, gin.H{"message":"User created successfully"})
}

func LoginHandler(c *gin.Context){
	userModel := db.User{}

	email := c.PostForm("email")
	password := c.PostForm("password")

	dbInstance := db.GetDB()

	if err := dbInstance.First(&userModel, "email = ?", email).Error; err != nil {
		c.JSON(401, gin.H{"error": "Invalid email or password"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(userModel.Password), []byte(password))
	if err != nil {
		c.JSON(401, gin.H{"error":"Invalid email or password"})
		return
	}
	session := sessions.Default(c)
	session.Set("userId", userModel.ID)
	session.Save()
	c.JSON(200, gin.H{"message": "Login successful", "user": userModel})
}

func CurrentUserHandler(c *gin.Context){
	userModel := db.User{}
	userID := sessions.Default(c).Get("userId")
	dbInstance := db.GetDB()

	if err := dbInstance.First(&userModel, userID).Error; err != nil {
		c.JSON(404, gin.H{"error":"User not found"})
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message":"Authorized",
		"user":userModel})
}

func LogoutHandler(c *gin.Context){
	session := sessions.Default(c)

	if session.Get("userId") != nil {
		session.Delete("userId")
		session.Save()
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message":"Logged out successfully"})
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil{
		return "", err
	}
	return string(hashedPassword), nil
}