package repositories

import (
	"context"
	"github.com/goravel/framework/facades"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"goravel/app/models"
	"goravel/app/transformers"
)

type CorporateInvoiceRepository struct {
	collection *mongo.Collection
}

func NewCorporateInvoiceRepository() *CorporateInvoiceRepository {
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
	return &CorporateInvoiceRepository{
		collection: client.Database(models.CorporateInvoice{}.DatabaseName()).Collection(models.CorporateInvoice{}.CollectionName()),
	}
}
func paginate(Page int64, TotalRecord int64, PageSize int64) (int64, int64) {
	totalPages := TotalRecord / PageSize
	if TotalRecord%PageSize != 0 {
		totalPages++
	}
	skip := (Page - 1) * PageSize
	return skip, totalPages
}
func (c CorporateInvoiceRepository) GetAll(Page int64, Limit int64, InvoiceNumber string) ([]*models.CorporateInvoice, *transformers.Pagination) {
	ctx := context.Background()
	var pipeline interface{}
	if InvoiceNumber == "" {
		pipeline = bson.M{}
	} else {
		pipeline = bson.D{{"invoiceNumber", bson.D{{"$regex", primitive.Regex{Pattern: InvoiceNumber}}}}}
	}
	TotalRecord, _ := c.collection.CountDocuments(ctx, pipeline)
	Skip, totalPages := paginate(Page, TotalRecord, int64(facades.Config().GetInt("app.pagination", 15)))
	cursor, err := c.collection.Find(ctx, pipeline, &options.FindOptions{
		Skip:  &Skip,
		Limit: &Limit,
	})
	if err != nil {
		facades.Log().Error(err)
	}
	var corporateInvoice []*models.CorporateInvoice
	if err := cursor.All(ctx, &corporateInvoice); err != nil {
		facades.Log().Error(err)
	}
	pagination := &transformers.Pagination{
		Total:       TotalRecord,
		CurrentPage: Page,
		TotalPages:  totalPages,
		PerPage:     Limit,
		Count:       int64(facades.Config().GetInt("app.pagination", 15)),
	}
	return corporateInvoice, pagination
}
