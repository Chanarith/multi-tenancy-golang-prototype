package middleware

import "github.com/gin-gonic/gin"

type Middleware interface {
	Handle(c *gin.Context)
}
