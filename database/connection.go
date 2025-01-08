package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/lamadev101/ecommerce-api/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type manager struct {
	connection *mongo.Client
	ctx        context.Context
	cancel     context.CancelFunc
}

var Mgr Manager

type Manager interface {
	Insert(interface{}, string) (interface{}, error)
	GetUserByEmail(string, string) *types.UserVerification
	UpateUserOTP(types.UserVerification, string) error
	UpdateEmailVerifiedStatus(types.UserVerification, string) error
	GetSingleRecordByEmailForUser(string, string) types.User
	UpdateByEmail(string, string, string) error
	GetListProducts(int, int, int, string) ([]types.Product, int64, error)
	CheckSlugOnDocument(string, string) error
	GetProductBySlug(string, string) (types.Product, error)
}

func ConnectDb() {
	uri := os.Getenv("MONGO_HOST")

	// Set client options
	clientOptions := options.Client().ApplyURI(uri)

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	fmt.Println("Connected to MongoDB successfully!")

	Mgr = &manager{connection: client, ctx: ctx, cancel: cancel}
}

// func DisconnectDb() {
// 	if Mgr != nil {
// 		if err := Mgr.connection.Disconnect(Mgr.ctx); err != nil {
// 			log.Fatalf("Failed to disconnect MongoDB: %v", err)
// 		}
// 		fmt.Println("Disconnected from MongoDB successfully!")
// 	}
// }

func Close(client *mongo.Client, ctx context.Context, cancel context.CancelFunc) {

	// CancelFunc to cancel to context
	defer cancel()

	// client provides a method to close
	// a mongoDB connection.
	defer func() {

		// client.Disconnect method also has deadline.
		// returns error if any,
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}
