package server

import (
	"Goose47/storage/internal/api/controllers"
	"Goose47/storage/internal/config"
	"fmt"
	"github.com/gin-gonic/gin"
)

// Serve runs http server and returns error if server stops
func Serve(
	cfg *config.Config,
	router *gin.Engine,
) error {
	gin.SetMode(cfg.Mode)

	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	return router.Run(addr)
}

func NewRouter(c *controllers.StorageController) *gin.Engine {
	r := gin.New()

	r.RedirectTrailingSlash = true
	r.RedirectFixedPath = true

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	AddApiGroup(r, c)

	return r
}
