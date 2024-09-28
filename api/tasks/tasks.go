package tasks

import (
	"Goose47/storage/config"
	database "Goose47/storage/db"
	"Goose47/storage/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"path"
	"time"
)

type TaskManager struct {
	stop chan bool
}

func (tm *TaskManager) RunTasks(
	db *database.DB,
	fsCfg *config.FSConfig,
) {
	go checkExpiredItems(tm.stop, db, fsCfg)
}

func (tm *TaskManager) StopTasks() {
	tm.stop <- true
	close(tm.stop)
}

func New() *TaskManager {
	tm := TaskManager{
		stop: make(chan bool, 1),
	}

	return &tm
}

func checkExpiredItems(
	stop <-chan bool,
	db *database.DB,
	cfg *config.FSConfig,
) {
	for {
		ticker := time.Tick(time.Second)
		select {
		case <-stop:
			return
		case <-ticker:
			cur, err := db.GetCollection().
				Find(
					context.TODO(),
					bson.D{
						{"exp", bson.D{{"$lte", time.Now().Unix()}}},
						{"exp", bson.D{{"$gt", 0}}},
					},
					options.Find().SetProjection(bson.D{{"_id", 1}, {"path", 1}}),
				)

			if err != nil {
				log.Println(err)
				continue
			}

			var results []models.StorageItem
			if err = cur.All(context.TODO(), &results); err != nil {
				log.Println(err)
				continue
			}

			for _, result := range results {
				nextPath := path.Join(cfg.Base, result.Path)
				if _, err := os.Stat(nextPath); os.IsNotExist(err) {
					continue
				}
				if err := os.Remove(nextPath); err != nil {
					log.Println(err)
					continue
				}
			}

			ids := make([]string, len(results))
			for i, r := range results {
				ids[i] = r.Key
			}

			_, err = db.GetCollection().
				DeleteMany(
					context.TODO(),
					bson.D{{"_id", bson.D{{"$in", ids}}}},
				)

			if err != nil {
				log.Println(err)
				continue
			}
		}
	}
}
