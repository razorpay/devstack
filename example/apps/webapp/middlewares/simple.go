package middlewares

import (
	"github.com/gin-gonic/gin"
	"log"
)

func SimpleMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("Logging from simple middleware")
		c.Next()
	}
}
