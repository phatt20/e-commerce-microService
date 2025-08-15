package userRepository

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"microService/config"
	"microService/modules/models"
	"microService/modules/payment"
	user "microService/modules/user"
	"microService/pkg/queue"
	"microService/pkg/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	UserRepository interface {
		GetOffset(pctx context.Context) (int64, error)
		UpserOffset(pctx context.Context, offset int64) error
		IsUniqueUser(pctx context.Context, email, username string) bool
		InsertOneUser(pctx context.Context, req *user.User) (primitive.ObjectID, error)
		FindOneUserProfile(pctx context.Context, userId string) (*user.UserProfileBson, error)
		InsertOneUserTranscation(pctx context.Context, req *user.UserTransaction) (primitive.ObjectID, error)
		GetUserSavingAccount(pctx context.Context, userId string) (*user.UserSavingAccount, error)
		FindOneUserCredential(pctx context.Context, email string) (*user.User, error)
		FindOneUserProfileToRefresh(pctx context.Context, userId string) (*user.User, error)
		DeleteOneUserTransaction(pctx context.Context, transactionId string) error
		DockedUserMoneyRes(pctx context.Context, cfg *config.Config, req *payment.PaymentTransferRes) error
		AddUserMoneyRes(pctx context.Context, cfg *config.Config, req *payment.PaymentTransferRes) error
	}
	userRepository struct {
		db *mongo.Client
	}
)

func NewUserRepository(db *mongo.Client) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) userDbConn(pctx context.Context) *mongo.Database {
	return r.db.Database("user")
}

func (r *userRepository) GetOffset(pctx context.Context) (int64, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.userDbConn(ctx)
	col := db.Collection("user_transactions_queue")

	result := new(models.KafkaOffset)
	if err := col.FindOne(ctx, bson.M{}).Decode(result); err != nil {
		log.Printf("Error: GetOffset failed: %s", err.Error())
		return -1, errors.New("error: GetOffset failed")
	}

	return result.Offset, nil
}

func (r *userRepository) UpserOffset(pctx context.Context, offset int64) error {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.userDbConn(ctx)
	col := db.Collection("user_transactions_queue")

	result, err := col.UpdateOne(ctx, bson.M{}, bson.M{"$set": bson.M{"offset": offset}}, options.Update().SetUpsert(true))
	if err != nil {
		log.Printf("Error: UpserOffset failed: %s", err.Error())
		return errors.New("error: UpserOffset failed")
	}
	log.Printf("Info: UpserOffset result: %v", result)

	return nil
}

func (r *userRepository) IsUniqueUser(pctx context.Context, email, username string) bool {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.userDbConn(ctx)
	col := db.Collection("users")

	user := new(user.User)
	if err := col.FindOne(
		ctx,
		bson.M{"$or": []bson.M{
			{"username": username},
			{"email": email},
		}},
	).Decode(user); err != nil {
		log.Printf("Error: IsUniqueUser: %s", err.Error())
		return true
	}
	return false
}

func (r *userRepository) InsertOneUser(pctx context.Context, req *user.User) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.userDbConn(ctx)
	col := db.Collection("users")

	userId, err := col.InsertOne(ctx, req)
	if err != nil {
		log.Printf("Error: InsertOneUser: %s", err.Error())
		return primitive.NilObjectID, errors.New("error: insert one user failed")
	}

	return userId.InsertedID.(primitive.ObjectID), nil
}

func (r *userRepository) DeleteOneUserTransaction(pctx context.Context, transactionId string) error {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.userDbConn(ctx)
	col := db.Collection("user_transactions")

	result, err := col.DeleteOne(ctx, bson.M{"_id": utils.ConvertToObjectId(transactionId)})
	if err != nil {
		log.Printf("Error: DeleteOneUserTransaction: %s", err.Error())
		return errors.New("error: delete one user transaction failed")
	}
	log.Printf("Delete result: %v", result)

	return nil
}

func (r *userRepository) FindOneUserProfile(pctx context.Context, userId string) (*user.UserProfileBson, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.userDbConn(ctx)
	col := db.Collection("users")

	result := new(user.UserProfileBson)

	if err := col.FindOne(
		ctx,
		bson.M{"_id": utils.ConvertToObjectId(userId)},
		options.FindOne().SetProjection(
			bson.M{
				"_id":        1,
				"email":      1,
				"username":   1,
				"created_at": 1,
				"updated_at": 1,
			},
		),
	).Decode(result); err != nil {
		log.Printf("Error: FindOneUserProfile: %s", err.Error())
		return nil, errors.New("error: user profile not found")
	}

	return result, nil
}

func (r *userRepository) InsertOneUserTranscation(pctx context.Context, req *user.UserTransaction) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.userDbConn(ctx)
	col := db.Collection("user_transactions")

	result, err := col.InsertOne(ctx, req)
	if err != nil {
		log.Printf("Error: InsertOneUserTranscation: %s", err.Error())
		return primitive.NilObjectID, errors.New("error: insert one user transcation failed")
	}
	log.Printf("Result: InsertOneUserTranscation: %v", result.InsertedID)

	return result.InsertedID.(primitive.ObjectID), nil
}

func (r *userRepository) GetUserSavingAccount(pctx context.Context, userId string) (*user.UserSavingAccount, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.userDbConn(ctx)
	col := db.Collection("user_transactions")

	filter := bson.A{
		bson.D{{"$match", bson.D{{"user_id", userId}}}},
		bson.D{
			{"$group",
				bson.D{
					{"_id", "$user_id"},
					{"balance", bson.D{{"$sum", "$amount"}}},
				},
			},
		},
		bson.D{
			{"$project",
				bson.D{
					{"user_id", "$_id"},
					{"_id", 0},
					{"balance", 1},
				},
			},
		},
	}

	cursors, err := col.Aggregate(ctx, filter)
	if err != nil {
		log.Printf("Error: GetUserSavingAccount: %s", err.Error())
		return nil, errors.New("error: failed to get user saving account")
	}

	result := new(user.UserSavingAccount)
	for cursors.Next(ctx) {
		if err := cursors.Decode(result); err != nil {
			log.Printf("Error: GetUserSavingAccount: %s", err.Error())
			return nil, errors.New("error: failed to get user saving account")
		}
	}

	return result, nil
}

func (r *userRepository) FindOneUserCredential(pctx context.Context, email string) (*user.User, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.userDbConn(ctx)
	col := db.Collection("users")

	result := new(user.User)

	if err := col.FindOne(ctx, bson.M{"email": email}).Decode(result); err != nil {
		log.Printf("Error: FindOneUserCredential: %s", err.Error())
		return nil, errors.New("error: email is invalid")
	}

	return result, nil
}

func (r *userRepository) FindOneUserProfileToRefresh(pctx context.Context, userId string) (*user.User, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.userDbConn(ctx)
	col := db.Collection("users")

	result := new(user.User)

	if err := col.FindOne(ctx, bson.M{"_id": utils.ConvertToObjectId(userId)}).Decode(result); err != nil {
		log.Printf("Error: FindOneUserProfileToRefresh: %s", err.Error())
		return nil, errors.New("error: user profile not found")
	}

	return result, nil
}

func (r *userRepository) DockedUserMoneyRes(pctx context.Context, cfg *config.Config, req *payment.PaymentTransferRes) error {
	reqInBytes, err := json.Marshal(req)
	if err != nil {
		log.Printf("Error: DockedUserMoneyRes failed: %s", err.Error())
		return errors.New("error: docked user money res failed")
	}

	if err := queue.PushMessageWithKeyToQueue(
		[]string{cfg.Kafka.Url},
		cfg.Kafka.ApiKey,
		cfg.Kafka.Secret,
		"payment",
		"buy",
		reqInBytes,
	); err != nil {
		log.Printf("Error: DockedUserMoneyRes failed: %s", err.Error())
		return errors.New("error: docked user money res failed")

	}

	return nil
}

func (r *userRepository) AddUserMoneyRes(pctx context.Context, cfg *config.Config, req *payment.PaymentTransferRes) error {
	reqInBytes, err := json.Marshal(req)
	if err != nil {
		log.Printf("Error: AddUserMoneyRes failed: %s", err.Error())
		return errors.New("error: docked user money res failed")
	}

	if err := queue.PushMessageWithKeyToQueue(
		[]string{cfg.Kafka.Url},
		cfg.Kafka.ApiKey,
		cfg.Kafka.Secret,
		"payment",
		"sell",
		reqInBytes,
	); err != nil {
		log.Printf("Error: AddUserMoneyRes failed: %s", err.Error())
		return errors.New("error: docked user money res failed")
	}

	return nil
}
