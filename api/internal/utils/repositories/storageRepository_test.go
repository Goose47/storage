package repositories

import (
	"Goose47/storage/internal/config"
	"Goose47/storage/internal/models"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

var mongoURI = "mongodb://localhost:27017"

func setupMongo(t *testing.T) (*mongo.Collection, func()) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	if err := client.Ping(context.TODO(), nil); err != nil {
		t.Fatalf("Failed to ping MongoDB: %v", err)
	}

	collection := client.Database("testdb").Collection("testcol")

	cleanup := func() {
		_ = collection.Drop(context.TODO())
		_ = client.Disconnect(context.TODO())
	}

	return collection, cleanup
}

func TestStorageRepository_FindByKey(t *testing.T) {
	if config.AppConfig.Mode == gin.ReleaseMode {
		t.Error("This tests can not be run in release mode")
	}

	coll, cleanup := setupMongo(t)
	defer cleanup()
	repo := NewStorageRepository(coll)

	key, item := "key", &models.StorageItem{"key", 0, "path", "name"}

	_, err := coll.InsertOne(context.TODO(), item)
	assert.NoError(t, err)

	_, err = repo.FindByKey("badKey")
	assert.Error(t, err)

	found, err := repo.FindByKey(key)
	assert.NoError(t, err)
	assert.Equal(t, *item, *found)

	//cleanup
	_, err = coll.DeleteOne(context.TODO(), bson.D{{"_id", key}})
	assert.NoError(t, err)
}

func TestStorageRepository_Set(t *testing.T) {
	if config.AppConfig.Mode == gin.ReleaseMode {
		t.Error("This tests can not be run in release mode")
	}

	coll, cleanup := setupMongo(t)
	defer cleanup()
	repo := NewStorageRepository(coll)

	key, item := "key", &models.StorageItem{"key", 0, "path", "name"}

	res, err := repo.Set(key, item)

	assert.NoError(t, err)
	assert.Equal(t, key, res)

	var foundItem models.StorageItem
	err = coll.
		FindOne(context.TODO(), bson.D{{"_id", key}}).
		Decode(&foundItem)

	assert.NoError(t, err)
	assert.Equal(t, *item, foundItem)

	//cleanup
	_, err = coll.DeleteOne(context.TODO(), bson.D{{"_id", key}})
	assert.NoError(t, err)
}

func TestStorageRepository_DeleteByKey(t *testing.T) {
	if config.AppConfig.Mode == gin.ReleaseMode {
		t.Error("This tests can not be run in release mode")
	}

	coll, cleanup := setupMongo(t)
	defer cleanup()
	repo := NewStorageRepository(coll)

	key, item := "key", &models.StorageItem{"key", 0, "path", "name"}

	_, err := repo.DeleteByKey(key)
	assert.Error(t, err)

	_, err = coll.InsertOne(context.TODO(), item)
	assert.NoError(t, err)

	deleted, err := repo.DeleteByKey(key)
	assert.NoError(t, err)
	assert.Equal(t, *item, *deleted)

	var foundItem models.StorageItem
	err = coll.
		FindOne(context.TODO(), bson.D{{"_id", key}}).
		Decode(&foundItem)
	assert.Error(t, err)
}
