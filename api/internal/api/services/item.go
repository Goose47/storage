package services

import (
	"Goose47/storage/internal/db"
	"Goose47/storage/internal/models"
	"Goose47/storage/internal/utils"
	"Goose47/storage/internal/utils/storage"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"log/slog"
	"mime/multipart"
	"path"
	"time"
)

type ItemSaver interface {
	SaveItem(key string, item *models.StorageItem) (string, error)
}

type ItemProvider interface {
	Item(key string) (*models.StorageItem, error)
}

type ItemDeleter interface {
	DeleteItem(key string) (*models.StorageItem, error)
}

type ItemService struct {
	log          *slog.Logger
	itemSaver    ItemSaver
	itemProvider ItemProvider
	itemDeleter  ItemDeleter
	storagePath  string
}

func NewItemService(
	log *slog.Logger,
	saver ItemSaver,
	provider ItemProvider,
	deleter ItemDeleter,
	storagePath string,
) *ItemService {
	return &ItemService{
		log:          log,
		itemSaver:    saver,
		itemProvider: provider,
		itemDeleter:  deleter,
		storagePath:  storagePath,
	}
}

func (s *ItemService) Get(key string) (*models.StorageItem, string, error) {
	const op = "api.services.item.Get"

	log := s.log.With(slog.String("op", op), slog.String("key", key))

	log.Info("attempting to retrieve item")

	item, err := s.itemProvider.Item(key)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Warn("item not found", slog.Any("error", err))
			return nil, "", fmt.Errorf("%s: %w", op, db.ErrItemNotFound)
		}

		log.Error("could not retrieve item", slog.Any("error", err))
		return nil, "", fmt.Errorf("%s: %w", op, err)
	}

	return item, path.Join(s.storagePath, item.Path), nil
}

func (s *ItemService) Set(
	key string,
	ttl int,
	file *multipart.FileHeader,
) (string, error) {
	const op = "api.services.item.Set"

	log := s.log.With(slog.String("op", op), slog.String("key", key))

	log.Info("searching for existing item")

	// Delete item and file if key exists
	existingItem, err := s.itemProvider.Item(key)
	if err == nil {
		log.Info("item exists, removing item")

		if err := storage.RemoveFileIfExists(path.Join(s.storagePath, existingItem.Path)); err != nil {
			log.Error("could not remove existing file", slog.Any("error", err))

			return "", fmt.Errorf("%s: %w", op, err)
		}
		if _, err := s.itemDeleter.DeleteItem(key); err != nil {
			log.Error("could not delete item", slog.Any("error", err))

			return "", fmt.Errorf("%s: %w", op, err)
		}

		log.Info("item removed")
	}

	log.Info("creating new item")

	item := &models.StorageItem{}
	// if exp == 0, document never expires
	exp := ttl
	if exp > 0 {
		exp += int(time.Now().Unix())
	}

	item.OriginalName = file.Filename
	item.Exp = exp
	item.Path = utils.GenerateRandomString(20) + path.Ext(file.Filename)

	log.Info("saving file")

	if err := storage.SaveFileFromHeader(file, path.Join(s.storagePath, item.Path)); err != nil {
		log.Error("could not save item", slog.Any("error", err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("file saved")
	log.Info("saving item")

	var id string
	if id, err = s.itemSaver.SaveItem(key, item); err != nil {
		log.Error("could not save item", slog.Any("error", err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("item saved")

	return id, nil
}

func (s *ItemService) Delete(key string) error {
	const op = "api.services.item.Delete"

	log := s.log.With(slog.String("op", op), slog.String("key", key))

	log.Info("deleting item")

	item, err := s.itemDeleter.DeleteItem(key)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Warn("item not found", slog.Any("error", err))
			return fmt.Errorf("%s: %w", op, db.ErrItemNotFound)
		}

		log.Error("could not delete item", slog.Any("error", err))

		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("item deleted")
	log.Info("removing file")

	err = storage.RemoveFileIfExists(path.Join(s.storagePath, item.Path))
	if err != nil {
		log.Error("could not remove file", slog.Any("error", err))

		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("file removed")

	return nil
}
