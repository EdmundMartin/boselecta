package mongoStorage

import (
	"context"
	"fmt"
	"github.com/EdmundMartin/boselecta/internal/flag"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoStorage struct {
	conn *mongo.Client
	coll *mongo.Collection
}

func createIndex(col *mongo.Collection) error {
	ctx, _ := context.WithTimeout(context.Background(), 10 * time.Second)
	mod := mongo.IndexModel{
		Keys: bson.M{
			"FlagName": -1, // index in descending order,
			"Namespace": -1,
		},
		Options: options.Index().SetUnique(true),
	}
	_, err := col.Indexes().CreateOne(ctx, mod)
	return err
}

func NewMongo(uri string, db string, collection string) (*MongoStorage, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}
	col := client.Database(db).Collection(collection)
	err = createIndex(col)
	if err != nil {
		return nil, err
	}
	return &MongoStorage{
		conn: client,
		coll: col,
	}, nil
}


func (ms *MongoStorage) All() ([]*flag.FeatureFlag, error) {
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cur, err := ms.coll.Find(ctx, bson.D{})
	if err != nil {
		return make([]*flag.FeatureFlag, 0), err
	}
	fflags := []*flag.FeatureFlag{}
	for cur.Next(ctx) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			return fflags, err
		}
		fflags = append(fflags, decodeBSON(result))
	}
	return fflags, nil
}


func (ms *MongoStorage) Create(namespace string, fl *flag.FeatureFlag) error {
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	_, err := ms.coll.InsertOne(ctx, bson.M{"Namespace": namespace, "FlagName": fl.FlagName, "Value": fl.Value,
		"Type": fl.Type.String(), "Refresh": fl.Refresh})
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func (ms *MongoStorage) GetFlag(namespace string, flagName string) (*flag.FeatureFlag, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10 * time.Second)
	var result bson.M
	err := ms.coll.FindOne(ctx, bson.M{"Namespace": namespace, "FlagName": flagName}).Decode(&result)
	if err != nil {
		return nil, err
	}
	return decodeBSON(result), nil
}

func (ms *MongoStorage) Delete(namespace string, flagName string) error {
	ctx, _ := context.WithTimeout(context.Background(), 10 * time.Second)
	_, err := ms.coll.DeleteOne(ctx, bson.M{"Namespace": namespace, "FlagName": flagName})
	return err
}

func (ms *MongoStorage) Update(namespace string, fl *flag.FeatureFlag) error {
	ctx, _ := context.WithTimeout(context.Background(), 10 * time.Second)
	_, err := ms.coll.UpdateOne(ctx, bson.M{"Namespace": namespace, "FlagName": fl.FlagName}, bson.M{
		"Namespace": namespace, "FlagName": fl.FlagName, "Value": fl.Value, "Type": fl.Type.String(),
		"Refresh": fl.Refresh,
	})
	return err
}