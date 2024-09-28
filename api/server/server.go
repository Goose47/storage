package server

import (
	"Goose47/storage/api/middleware"
	"Goose47/storage/config"
	"fmt"
	"github.com/gin-gonic/gin"
)

func Serve(cfg *config.Config) error {
	r := NewRouter()
	gin.SetMode(cfg.Mode)

	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	return r.Run(addr)
}

func NewRouter() *gin.Engine {
	r := gin.New()

	r.RedirectTrailingSlash = true
	r.RedirectFixedPath = true

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.Handle404)

	AddApiGroup(r)

	return r
}
