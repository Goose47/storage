package db

import (
	"Goose47/storage/config"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

type DB struct {
	Conn *mongo.Client
	cfg  *config.DBConfig
}

func (db *DB) GetCollection() *mongo.Collection {
	return db.Conn.Database(db.cfg.DBName).Collection(db.cfg.DBColl)
}

func New(cfg *config.DBConfig) (*DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	conn, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Url))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %s", cfg.Url)
	}
	err = conn.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, fmt.Errorf("failed to ping database")
	}

	return &DB{
		Conn: conn,
		cfg:  cfg,
	}, nil
}
