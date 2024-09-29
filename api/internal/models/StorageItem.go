package models

type StorageItem struct {
	Key          string `bson:"_id"`
	Exp          int    `bson:"exp"`
	Path         string `bson:"path"`
	OriginalName string `bson:"originalName"`
}
