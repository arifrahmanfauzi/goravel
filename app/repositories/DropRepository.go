package repositories

import (
	"context"
	"github.com/goravel/framework/facades"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"goravel/app/models"
)

type DropRepository struct {
	collection *mongo.Collection
}

func NewDropRepository() *DropRepository {
	clientOptions := options.Client().ApplyURI(facades.Config().GetString("DB_STRING", ""))
	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		facades.Log().Error(err)
	}
	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		facades.Log().Error(err)
	}
	drop := models.Drop{}
	return &DropRepository{
		collection: client.Database(drop.DatabaseName()).Collection(drop.CollectionName()),
	}
}

func (Drop *DropRepository) GetAll(page int64, pageSize int64, total *int64, totalPage *int64) ([]*models.Drop, error) {
	ctx := context.Background()
	// Calculate the number of documents to skip
	skip := (page - 1) * pageSize
	// Get total count of trips
	totalRecord, err := Drop.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	*total = totalRecord
	// Calculate total pages
	totalPages := totalRecord / pageSize
	if totalRecord%pageSize != 0 {
		totalPages++
	}
	*totalPage = totalPages

	cursor, err := Drop.collection.Find(ctx, bson.M{}, &options.FindOptions{
		Skip:  &skip,
		Limit: &pageSize,
	})
	if err != nil {
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {

		}
	}(cursor, ctx)

	var Drops []*models.Drop
	for cursor.Next(ctx) {
		var customerOrder models.Drop
		if err := cursor.Decode(&customerOrder); err != nil {
			return nil, err
		}
		Drops = append(Drops, &customerOrder)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return Drops, nil
}

func (Drop *DropRepository) Update(ID string, update map[string]interface{}) *models.Drop {
	id, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		facades.Log().Error(err)
		return nil
	}
	filter := bson.M{"_id": id}
	update = bson.M{"$set": update}
	// Set returnDocument for returning the updated document
	returnDocument := options.FindOneAndUpdate().SetReturnDocument(options.After)
	//Update content
	var Drops *models.Drop
	err = Drop.collection.FindOneAndUpdate(context.Background(), filter, update, returnDocument).Decode(&Drops)
	if err != nil {
		facades.Log().Error(err)
		return nil
	}

	return Drops
}

func (Drop *DropRepository) Delete(ID string) (*mongo.DeleteResult, error) {
	id, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		facades.Log().Error(err)
		return nil, err
	}
	filter := bson.M{"_id": id}
	res, err := Drop.collection.DeleteOne(context.Background(), filter)
	if err != nil {
		facades.Log().Error(err)
	}
	return res, err
}
