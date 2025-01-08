package database

import (
	"context"
	"fmt"
	"log"

	"github.com/lamadev101/ecommerce-api/constant"
	"github.com/lamadev101/ecommerce-api/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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

// ================ Products ====================
func (mgr *manager) GetListProducts(page, limit, offset int, collectionName string) (products []types.Product, count int64, err error) {
	skip := ((page - 1) * limit)
	if offset > 0 {
		skip = offset
	}

	orgCollection := mgr.connection.Database(constant.DATABASE).Collection(collectionName)

	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))

	cur, err := orgCollection.Find(context.TODO(), bson.M{}, findOptions)
	if err != nil {
		return nil, 0, err
	}
	if err = cur.All(context.TODO(), &products); err != nil {
		return nil, 0, err
	}

	itemCount, err := orgCollection.CountDocuments(context.TODO(), bson.M{})
	return products, itemCount, err
}

func (mgr *manager) CheckSlugOnDocument(slug string, collectionName string) error {
	res := &types.Product{}
	filter := bson.D{{Key: "slug", Value: slug}}
	orgCollection := mgr.connection.Database(constant.DATABASE).Collection(collectionName)

	err := orgCollection.FindOne(context.TODO(), filter).Decode(&res)
	return err
}

func (mgr *manager) GetProductBySlug(slug, collectionName string) (product types.Product, err error) {
	orgCollection := mgr.connection.Database(constant.DATABASE).Collection(collectionName)
	filter := bson.M{"slug": slug}

	err = orgCollection.FindOne(context.TODO(), filter).Decode(&product)
	if err != nil {
		log.Println("error:", err)
		return product, nil
	}

	return product, nil
}
