package datasources

import (
	"context"
	"log"
	"os"

	"go.elastic.co/apm/module/apmmongo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	Context context.Context
	MongoDB *mongo.Client
}

func NewMongoDB(maxPoolSize uint64) *MongoDB {

	option := options.Client().ApplyURI(os.Getenv("MONGODB_URI")).SetMonitor(apmmongo.CommandMonitor()).SetMaxPoolSize(maxPoolSize)
	client, err0 := mongo.Connect(context.Background(), option)

	if err0 != nil {
		log.Fatal("error connection : ", err0)
	}

	return &MongoDB{
		Context: context.Background(),
		MongoDB: client,
	}
}
