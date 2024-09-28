package services

import (
	"Goose47/storage/db"
	"Goose47/storage/models"
	"Goose47/storage/utils"
	"Goose47/storage/utils/repositories"
	"Goose47/storage/utils/storage"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"path"
	"time"
)

type ItemService struct {
	repo *repositories.StorageRepository
}

func NewItemService() *ItemService {
	return &ItemService{repositories.NewStorageRepository(db.GetCollection())}
}

func (service *ItemService) Set(
	key string,
	ttl int,
	file *multipart.FileHeader,
	c *gin.Context,
) (string, error) {
	// Delete item and file if key exists
	existingItem, err := service.repo.FindByKey(key)
	if err == nil {
		if err := storage.RemoveFileIfExists(existingItem.GetFullPath()); err != nil {
			return "", err
		}
		if _, err := service.repo.DeleteByKey(key); err != nil {
			return "", err
		}
	}

	item := &models.StorageItem{}
	// if exp == 0, document never expires
	exp := ttl
	if exp > 0 {
		exp += int(time.Now().Unix())
	}

	item.OriginalName = file.Filename
	item.Exp = exp
	item.Path = utils.GenerateRandomString(20) + path.Ext(file.Filename)

	if err := c.SaveUploadedFile(file, item.GetFullPath()); err != nil {
		return "", err
	}

	var id string
	if id, err = service.repo.Set(key, item); err != nil {
		return "", err
	}

	return id, nil
}

func (service *ItemService) Delete(key string) error {
	item, err := service.repo.DeleteByKey(key)

	if err != nil {
		return err
	}

	err = storage.RemoveFileIfExists(item.GetFullPath())
	if err != nil {
		return err
	}

	return nil
}
