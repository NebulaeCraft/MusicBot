package bff

import "github.com/gin-gonic/gin"

func AddNewHandler(s *gin.Engine) {
	s.GET("/test", func(c *gin.Context) {
		c.String(200, "Hello World!")
	})
}
