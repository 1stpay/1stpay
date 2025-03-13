package test

import (
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

func PerformRequest(router *gin.Engine, method, url string, body io.Reader, bearerToken ...string) *httptest.ResponseRecorder {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	if len(bearerToken) > 0 && bearerToken[0] != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken[0])
	}
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	return recorder
}
