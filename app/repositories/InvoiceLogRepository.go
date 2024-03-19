package repositories

import (
	"context"
	"fmt"
	"github.com/goravel/framework/facades"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"goravel/app/models"
	"time"
)

type InvoiceLogRepository struct {
	collection *mongo.Collection
}

func NewInvoiceLogRepository() *InvoiceLogRepository {
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
	invoiceLog := models.InvoiceLog{}
	return &InvoiceLogRepository{
		collection: client.Database(invoiceLog.DatabaseName()).Collection(invoiceLog.CollectionName()),
	}
}

func (i *InvoiceLogRepository) Create(model interface{}) error {
	//ctx := context.Background()
	//var InvoiceLogs *models.InvoiceLog
	//i.collection.InsertOne(ctx, model)
	return nil
}
func (i *InvoiceLogRepository) GetAll(Page int64, PageSize int64) ([]*models.InvoiceLog, int64, int64, error) {
	ctx := context.Background()
	// Calculate the number of documents to skip
	skip := (Page - 1) * PageSize
	// Get total count of drops
	totalRecord, err := i.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, totalRecord, 0, err
	}
	// Calculate total pages
	totalPages := totalRecord / PageSize
	if totalRecord%PageSize != 0 {
		totalPages++
	}

	// Set total and totalPage
	//total := totalRecord
	totalPage := totalPages
	// Find drops
	cursor, err := i.collection.Find(ctx, bson.M{}, &options.FindOptions{
		Skip:  &skip,
		Limit: &PageSize,
	})
	if err != nil {
		return nil, totalRecord, 0, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			facades.Log().Error(err)
		}
	}(cursor, ctx)
	var InvoiceLog []*models.InvoiceLog
	if err := cursor.All(context.Background(), &InvoiceLog); err != nil {
		return nil, totalRecord, 0, err
	}
	return InvoiceLog, totalRecord, totalPage, nil
}

func (i *InvoiceLogRepository) FindById(Id string) (*models.InvoiceLog, error) {
	var invoiceLog models.InvoiceLog
	ID, err := primitive.ObjectIDFromHex(Id)
	err = i.collection.FindOne(context.Background(), bson.M{"_id": ID}).Decode(&invoiceLog)
	if err != nil {
		return nil, err
	}
	return &invoiceLog, err
}

func (i *InvoiceLogRepository) FindByField(value string) []*models.InvoiceLog {
	filter := bson.M{
		"invoiceNumber": value,
	}
	//var invoiceLog []models.InvoiceLog
	cursor, err := i.collection.Find(context.Background(), filter)
	if err != nil {
		facades.Log().Error(err)
		return nil
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			facades.Log().Error(err)
		}
	}(cursor, context.Background())
	var invoiceLogs []*models.InvoiceLog
	for cursor.Next(context.Background()) {
		var invoiceLog models.InvoiceLog
		if err := cursor.Decode(&invoiceLog); err != nil {
			fmt.Println("Error decoding document:", err)
			continue
		}
		invoiceLogs = append(invoiceLogs, &invoiceLog)
	}
	return invoiceLogs
}

func (i *InvoiceLogRepository) Delete(Id string) (*mongo.DeleteResult, error) {
	id, err := primitive.ObjectIDFromHex(Id)
	if err != nil {
		facades.Log().Error(err)
	}
	filter := bson.M{"_id": id}
	res, err := i.collection.DeleteOne(context.Background(), filter)
	if err != nil {
		facades.Log().Error(err)
	}
	return res, err
}

func (i *InvoiceLogRepository) Update(Id string, update map[string]interface{}) *models.InvoiceLog {
	id, err := primitive.ObjectIDFromHex(Id)
	if err != nil {
		facades.Log().Error(err)
		return nil
	}
	filter := bson.M{"_id": id}
	update["updated_at"] = time.Now()
	update = bson.M{"$set": update}
	// Set returnDocument for returning the updated document
	returnDocument := options.FindOneAndUpdate().SetReturnDocument(options.After)
	//update content
	var InvoiceLogs *models.InvoiceLog
	err = i.collection.FindOneAndUpdate(context.Background(), filter, update, returnDocument).Decode(&InvoiceLogs)
	if err != nil {
		facades.Log().Error(err)
		return nil
	}
	return InvoiceLogs
}
