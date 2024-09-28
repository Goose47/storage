package models

import (
	"Goose47/storage/config"
	"path"
)

type StorageItem struct {
	Key          string `bson:"_id"`
	Exp          int    `bson:"exp"`
	Path         string `bson:"path"`
	OriginalName string `bson:"originalName"`
}

func (i *StorageItem) GetFullPath() string {
	return path.Join(config.FSConfig.Base, i.Path)
}
