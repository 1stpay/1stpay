package config

import "github.com/gin-gonic/gin"

type Middleware struct {
	JWTAuth    gin.HandlerFunc
	APIKeyAuth gin.HandlerFunc
}
