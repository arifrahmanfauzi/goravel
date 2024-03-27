package repositories

import (
	"context"
	"github.com/goravel/framework/facades"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"goravel/app/models"
	"time"
)

type TripRepository struct {
	collection *mongo.Collection
}

func NewTripRepository() *TripRepository {
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
	facades.Log().Info("mongoDB has connected!")
	trip := models.Trip{}
	return &TripRepository{
		collection: client.Database(trip.DatabaseName()).Collection(trip.CollectionName()),
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
	facades.Log().Info(totalRecord)
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
func (tr *TripRepository) FindTripByOrderNumber() []*models.Trip {
	ctx := context.Background()
	//lookupStage := bson.D{
	//	{"$lookup",
	//		bson.D{
	//			{"from", "customerOrder"},
	//			{"localField", "orderIdObject"},
	//			{"foreignField", "_id"},
	//			{"as", "cO"},
	//		},
	//	},
	//}
	pipeLine := bson.A{
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "customerOrder"},
					{"localField", "orderIdObject"},
					{"foreignField", "_id"},
					{"as", "cO"},
				},
			},
		},
		bson.D{
			{"$unwind",
				bson.D{
					{"path", "$cO"},
					{"preserveNullAndEmptyArrays", false},
				},
			},
		},
		bson.D{{"$match", bson.D{{"cO.orderNumber", bson.D{{"$eq", "SD-2023070525b5f8"}}}}}},
		bson.D{
			{"$project",
				bson.D{
					{"_id", 1},
					{"orderId", 1},
					{"tripNumber", 1},
					{"distance", 1},
					{"pick", 1},
				},
			},
		},
	}
	cursor, err := tr.collection.Aggregate(ctx, pipeLine, options.Aggregate().SetMaxTime(5*time.Second))
	//cursor, err := tr.collection.Find(ctx, bson.D{})
	if err != nil {
		facades.Log().Error(err)
		return nil
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			facades.Log().Error(err)
		}
	}(cursor, ctx)
	var trips []*models.Trip
	if err := cursor.All(ctx, &trips); err != nil {
		facades.Log().Error(err)
	}
	if err != nil {
		facades.Log().Error(err)
		return nil
	}
	return trips
}
