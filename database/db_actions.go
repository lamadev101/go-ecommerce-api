package database

import (
	"context"
	"fmt"

	"github.com/lamadev101/ecommerce-api/constant"
	"github.com/lamadev101/ecommerce-api/types"
	"go.mongodb.org/mongo-driver/bson"
)

func (mgr *manager) Insert(document interface{}, collectionName string) (interface{}, error) {
	fmt.Println("document: ", document)

	orgCollection := mgr.connection.Database(constant.DATABASE).Collection(collectionName)

	// insert the bson object using InsertOne() method
	result, err := orgCollection.InsertOne(context.TODO(), document)
	fmt.Println("result: ", result, err)

	// check for error in the insertion
	if err != nil {
		return nil, err
	}
	return result.InsertedID, nil
}

func (mgr *manager) GetUserByEmail(email string, collectionName string) *types.UserVerification {
	res := &types.UserVerification{}
	filter := bson.D{{Key: "email", Value: email}}
	orgCollection := mgr.connection.Database(constant.DATABASE).Collection(collectionName)

	_ = orgCollection.FindOne(context.TODO(), filter).Decode(&res)
	return res
}

func (mgr *manager) UpateUserOTP(data types.UserVerification, collectionName string) error {
	orgCollection := mgr.connection.Database(constant.DATABASE).Collection(collectionName)

	filter := bson.D{{Key: "email", Value: data.Email}}
	update := bson.D{{Key: "$set", Value: data}}

	_, err := orgCollection.UpdateOne(context.TODO(), filter, update)

	return err
}

func (mgr *manager) UpdateEmailVerifiedStatus(req types.UserVerification, collectionName string) error {
	orgCollection := mgr.connection.Database(constant.DATABASE).Collection(collectionName)
	filter := bson.D{{Key: "email", Value: req.Email}}
	update := bson.D{{Key: "$set", Value: req}}

	_, err := orgCollection.UpdateOne(context.TODO(), filter, update)
	return err
}

func (mgr *manager) UpdateByEmail(email string, hashPassword string, collectionName string) error {
	orgCollection := mgr.connection.Database(constant.DATABASE).Collection(collectionName)
	filter := bson.D{{Key: "email", Value: email}}
	update := bson.M{"$set": bson.M{"password": hashPassword}}

	_, err := orgCollection.UpdateOne(context.TODO(), filter, update)
	return err
}

func (mgr *manager) GetSingleRecordByEmailForUser(email, collectionName string) types.User {
	res := types.User{}
	filter := bson.D{{Key: "email", Value: email}}

	orgCollection := mgr.connection.Database(constant.DATABASE).Collection(collectionName)

	_ = orgCollection.FindOne(context.TODO(), filter).Decode(&res)
	return res
}
