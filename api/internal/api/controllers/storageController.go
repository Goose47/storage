package controllers

import (
	"Goose47/storage/internal/api/services"
	"Goose47/storage/internal/db"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"net/http"
)

type StorageController struct {
	service *services.ItemService
}

func NewStorageController(
	service *services.ItemService,
) *StorageController {
	return &StorageController{
		service: service,
	}
}

func (con *StorageController) Get(c *gin.Context) {
	key := c.Param("key")

	item, downloadPath, err := con.service.Get(key)

	if err != nil {
		if errors.Is(err, db.ErrItemNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("%s is not found", key)})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
		return
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+item.OriginalName)
	c.File(downloadPath)
}

type SetForm struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
	Ttl  int                   `form:"ttl"`
}

func (con *StorageController) Set(c *gin.Context) {
	key := c.Param("key")

	var form SetForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	id, err := con.service.Set(key, form.Ttl, form.File)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
		return
	}

	c.JSON(200, gin.H{
		"message": fmt.Sprintf("Set %s", id),
	})
}

func (con *StorageController) Delete(c *gin.Context) {
	key := c.Param("key")

	err := con.service.Delete(key)

	if err != nil {
		if errors.Is(err, db.ErrItemNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("%s is not found", key)})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
		return
	}

	c.JSON(200, gin.H{
		"message": "ok",
	})
}
