package server

import (
	"Goose47/storage/internal/api/controllers"
	"github.com/gin-gonic/gin"
)

func AddApiGroup(
	r *gin.Engine,
	c *controllers.StorageController,
) {
	api := r.Group("api")
	{
		v1 := api.Group("v1")
		{
			v1.GET("healthcheck")

			storage := v1.Group("storage")
			{
				storage.GET("/:key", c.Get)
				storage.POST("/:key", c.Set)
				storage.DELETE("/:key", c.Delete)
			}
		}
	}
}
