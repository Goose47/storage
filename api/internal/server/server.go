package server

import (
	"Goose47/storage/internal/api/controllers"
	"Goose47/storage/internal/config"
	"Goose47/storage/internal/server/middleware"
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
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

func NewRouter(
	log *slog.Logger,
	c *controllers.StorageController,
	secret string,
	permsProvider middleware.PermsProvider,
) *gin.Engine {
	r := gin.New()

	r.RedirectTrailingSlash = true
	r.RedirectFixedPath = true

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.NewAuthMiddleware(log, secret, permsProvider))

	AddApiGroup(r, c)

	return r
}
