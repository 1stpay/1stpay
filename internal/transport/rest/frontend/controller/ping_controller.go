package controller

import "github.com/gin-gonic/gin"

type PingSuccess struct {
	Status  int  `json:"status"`
	Success bool `json:"success"`
}

func Ping(c *gin.Context) {
	responseData := PingSuccess{Status: 200, Success: true}
	c.JSON(200, responseData)
}
