package database

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	JobsCollection  *mongo.Collection
	UsersCollection *mongo.Collection
	Database        *mongo.Database
	Client          *mongo.Client
}

func (db *Database) InitDatabase() error {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv("DB_URI")))
	if err != nil {
		return err
	}
	db.Client = client
	db.Database = db.Client.Database(os.Getenv("DB_DATABASE"))
	db.JobsCollection = db.Database.Collection(os.Getenv("DB_JOBS_COLLECTION"))
	db.UsersCollection = db.Database.Collection(os.Getenv("DB_USER_COLLECTION"))
	return nil
}

var Db = &Database{}

func (db *Database) CheckUsernameAvailability(username string) (bool, error) {
	cur, err := Db.UsersCollection.Find(context.Background(), bson.M{"username": username})
	if err != nil {
		return false, err
	}
	if cur.RemainingBatchLength() != 0 {
		return false, nil
	}
	return true, nil
}

func (db *Database) GetUser(username string) (interface{}, error) {
	var user interface{}
	err := Db.UsersCollection.FindOne(context.Background(), bson.M{"username": username}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (db *Database) UpdateUser(newUser interface{}) error {
	userByte, err := json.Marshal(newUser)
	if err != nil {
		return err
	}
	var user map[string]interface{}
	err = json.Unmarshal(userByte, &user)
	if err != nil {
		return err
	}
	res, err := Db.UsersCollection.UpdateOne(context.Background(), bson.M{"_id": user["id"].(string)}, bson.M{
		"$set": newUser,
	})
	if err != nil {
		return err
	}
	if res.ModifiedCount == 1 {
		return nil
	}
	return fmt.Errorf("update user failed, try again later")
}
