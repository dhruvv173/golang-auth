package handlers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func IsAuth() gin.HandlerFunc{
	return func (c *gin.Context)  {
		session := sessions.Default(c)

		userID := session.Get("userId")

		if userID == nil{
			c.JSON(http.StatusUnauthorized, gin.H{"error":"Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}