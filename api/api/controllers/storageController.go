package controllers

import (
	"Goose47/storage/api/errs"
	"Goose47/storage/api/services"
	"Goose47/storage/db"
	"Goose47/storage/utils/repositories"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"mime/multipart"
	"net/http"
)

type StorageController struct {
	repo *repositories.StorageRepository
}

func NewStorageController() *StorageController {
	return &StorageController{repositories.NewStorageRepository(db.GetCollection())}
}

func (con *StorageController) Get(c *gin.Context) {
	key := c.Param("key")

	item, err := con.repo.FindByKey(key)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.Error(&errs.NotFoundError{Message: fmt.Sprintf("%s is not found", key)})
			return
		}
		log.Panic(err)
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+item.OriginalName)
	c.File(item.GetFullPath())
}

type SetForm struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
	Ttl  int                   `form:"ttl"`
}

func (*StorageController) Set(c *gin.Context) {
	itemService := services.NewItemService()
	key := c.Param("key")

	var form SetForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	id, err := itemService.Set(key, form.Ttl, form.File, c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": fmt.Sprintf("Set %s", id),
	})
}

func (*StorageController) Delete(c *gin.Context) {
	itemService := services.NewItemService()
	key := c.Param("key")

	err := itemService.Delete(key)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.Error(&errs.NotFoundError{Message: fmt.Sprintf("%s is not found", key)})
			return
		}
		log.Panic(err)
	}

	c.JSON(200, gin.H{
		"message": "ok",
	})
}
