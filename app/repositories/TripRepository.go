package repositories

import (
	"context"
	"goravel/app/models"

	"github.com/goravel/framework/facades"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TripRepository struct {
	collection *mongo.Collection
}

func NewTripRepository(client *mongo.Client, dbName, collectionName string) *TripRepository {
	return &TripRepository{
		collection: client.Database(dbName).Collection(collectionName),
	}
}

func (tr *TripRepository) Create(trip *models.Trip) error {
	_, err := tr.collection.InsertOne(context.Background(), trip)
	return err
}

func (tr *TripRepository) FindByID(id string) (*models.Trip, error) {
	var trip models.Trip
	err := tr.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&trip)
	if err != nil {
		return nil, err
	}
	return &trip, nil
}
func (tr *TripRepository) GetAll(page, pageSize int64, total *int64, totalPage *int64) ([]*models.Trip, error) {
	ctx := context.Background()

	// Calculate the number of documents to skip
	skip := (page - 1) * pageSize

	// Get total count of trips
	totalRecord, err := tr.collection.CountDocuments(ctx, bson.M{})
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
	// Find the trips
	cursor, err := tr.collection.Find(ctx, bson.M{}, &options.FindOptions{
		Skip:  &skip,
		Limit: &pageSize,
	})
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			// Handle or log the error if closing the cursor fails
			facades.Log().Error(err)
		}
	}()

	var trips []*models.Trip
	for cursor.Next(ctx) {
		var trip models.Trip
		if err := cursor.Decode(&trip); err != nil {
			return nil, err
		}
		trips = append(trips, &trip)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return trips, nil
}
