package repositories

import (
	"Goose47/storage/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type StorageRepository struct {
	collection *mongo.Collection
}

func NewStorageRepository(collection *mongo.Collection) *StorageRepository {
	return &StorageRepository{collection}
}

func (repo *StorageRepository) FindByKey(key string) (*models.StorageItem, error) {
	var result bson.M
	err := repo.collection.
		FindOne(context.TODO(), bson.D{{"_id", key}}).
		Decode(&result)

	if err != nil {
		return nil, err
	}

	return createItem(result), nil
}

func (repo *StorageRepository) Set(key string, item *models.StorageItem) (string, error) {
	result, err := repo.collection.
		InsertOne(context.TODO(), bson.D{
			{"_id", key},
			{"path", item.Path},
			{"exp", item.Exp},
			{"originalName", item.OriginalName},
		})

	if err != nil {
		return "", err
	}

	return result.InsertedID.(string), nil
}

func (repo *StorageRepository) DeleteByKey(key string) (*models.StorageItem, error) {
	var result bson.M
	err := repo.collection.
		FindOneAndDelete(context.TODO(), bson.D{{"_id", key}}).
		Decode(&result)

	if err != nil {
		return nil, err
	}

	return createItem(result), nil
}

func createItem(result map[string]interface{}) *models.StorageItem {
	item := &models.StorageItem{}

	item.Key = result["_id"].(string)
	item.Path = result["path"].(string)
	item.Exp = int(result["exp"].(int32))
	item.OriginalName = result["originalName"].(string)

	return item
}
