package repositories

import (
	"context"
	"github.com/goravel/framework/facades"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"goravel/app/models"
)

type OrderRepository struct {
	collection *mongo.Collection
}

func NewOrderRepository() *OrderRepository {
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
	return &OrderRepository{
		collection: client.Database(models.CustomerOrder{}.DatabaseName()).Collection(models.CustomerOrder{}.CollectionName()),
	}
}

func (or *OrderRepository) GetAll(page, pageSize int64, total *int64, totalPage *int64) ([]*models.CustomerOrder, error) {
	ctx := context.Background()

	// Calculate the number of documents to skip
	skip := (page - 1) * pageSize

	// Get total count of trips
	totalRecord, err := or.collection.CountDocuments(ctx, bson.M{})
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

	cursor, err := or.collection.Find(ctx, bson.M{}, &options.FindOptions{
		Skip:  &skip,
		Limit: &pageSize,
	})
	if err != nil {
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			facades.Log().Error(err)
		}
	}(cursor, ctx)

	var CO []*models.CustomerOrder
	for cursor.Next(ctx) {
		var customerOrder models.CustomerOrder
		if err := cursor.Decode(&customerOrder); err != nil {
			return nil, err
		}
		CO = append(CO, &customerOrder)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return CO, nil
}

func (OrderRepository) FindById() {
	//TODO implement me
	panic("implement me")
}

func (OrderRepository) Create() {
	//TODO implement me
	panic("implement me")
}

func (OrderRepository) Update() {
	//TODO implement me
	panic("implement me")
}

func (OrderRepository) Delete() {
	//TODO implement me
	panic("implement me")
}
