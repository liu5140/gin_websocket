package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	hmacSampleSecret = "tnb9Y0du$2a$10$KmatydruRTKlaUwErUOtNOXiPHVPunb9Y0dup9newm"
	hmacSecureSecret = "Tnb9Y0du$2a$10$KmatydruRTKlaUwErUOtNOXiPHVPunb9Y0dup9newm"
)

func MainRouter() http.Handler {
	engine := gin.New()
	router := engine.Group("acc")
	webSocketRouter(router)
	return engine
}
