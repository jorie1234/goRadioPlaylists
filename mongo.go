package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"log"
	"time"
)

type mongodb struct {
	client *mongo.Client
}

func NewMongoDB() (*mongodb, error) {

	//username := os.Getenv("MONGODB_USERNAME")
	//password := os.Getenv("MONGODB_PASSWORD")
	//clusterEndpoint := os.Getenv("MONGODB_ENDPOINT")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error
	m := mongodb{}
	m.client, err = mongo.Connect(ctx, options.Client().ApplyURI(
		"mongodb://mongo:27017/radio?retryWrites=true&w=majority",
	))
	if err != nil {
		log.Fatal(err)
	}

	// Force a connection to verify our connection string
	err = m.client.Ping(ctx, nil)
	if err != nil {
		log.Printf("Failed to ping cluster: %v", err)
	}
	return &m, nil
}
func (m *mongodb) Insert(table string, data interface{}) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := m.client.Database("radio").Collection(table).InsertOne(ctx, data)
	if err != nil {
		log.Printf("Could not create Entry: %v", err)
		return
	}
	oid := result.InsertedID.(primitive.ObjectID)
	log.Printf("Successfully inserted data oid %v", oid)
	return

}

func (m *mongodb) GetLastEntry(table string) PlaylistEntry {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var data []PlaylistEntry

	options := options.Find()

	// Sort by `created` field descending
	options.SetSort(bson.D{{"created", -1}})

	// Limit by 10 documents only
	options.SetLimit(1)

	cursor, err := m.client.Database("radio").Collection(table).Find(ctx, bson.D{}, options)
	if err != nil {
		log.Printf("find failed: %v", err)
		return PlaylistEntry{}
	}
	cursor.All(ctx, &data)
	if len(data) == 1 {
		return data[0]
	}
	return PlaylistEntry{}

}
func (m *mongodb) EnsureIndex(table string) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := m.client.Database("radio").Collection(table).Indexes().CreateOne(
		ctx,
		mongo.IndexModel{
			Keys:    bsonx.Doc{{"created", bsonx.Int32(1)}},
			Options: options.Index().SetUnique(true),
		},
	)

	if err != nil {
		log.Printf("create index on table %s  failed: %v", table, err)
		return
	}
	return
}
