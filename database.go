package main

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
	db.UsersCollection = db.Database.Collection(os.Getenv("DB_USER_COLLECTION"))
	return nil
}

var Db = &Database{}

func (db *Database) CheckUsernameAvailability(username string) (bool, error) {
	cur, err := db.UsersCollection.Find(context.Background(), bson.M{"username": username})
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
	err := db.UsersCollection.FindOne(context.Background(), bson.M{"username": username}).Decode(&user)
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
	var updatedUser interface{}
	err = db.UsersCollection.FindOneAndReplace(context.Background(), bson.M{"_id": user["id"].(string)}, newUser, nil).Decode(&updatedUser)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) GetJob(username string, jobId string) (interface{}, error) {
	var user interface{}
	// var job interface{}
	err := db.UsersCollection.FindOne(context.Background(), bson.M{"username": username, "jobs": bson.M{"$elemMatch": bson.M{"_id": jobId}}}).Decode(&user)
	if err != nil {
		return nil, err
	}
	fmt.Println(user)
	return nil, nil
}

func (db *Database) UpdateJob(jobId string, newUser interface{}) error {
	userByte, err := bson.Marshal(newUser)
	if err != nil {
		return err
	}
	var userBson map[string]interface{}
	err = bson.Unmarshal(userByte, &userBson)
	if err != nil {
		return err
	}
	filter := bson.D{{"_id", userBson["_id"].(string)}, {"jobs._id", jobId}}
	var user interface{}

	fmt.Println("jobid", jobId, userBson["nuclei_result"])
	err = db.UsersCollection.FindOneAndReplace(context.TODO(), filter, newUser, nil).Decode(&user)
	if err != nil {
		return err
	}
	return nil
}
