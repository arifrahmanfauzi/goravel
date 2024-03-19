package models

import (
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type InvoiceLog struct {
	ID            string    `bson:"_id,omitempty"`
	InvoiceId     string    `bson:"invoiceId,omitempty"`
	InvoiceNumber string    `bson:"invoiceNumber,omitempty"`
	AdminEmail    string    `bson:"adminEmail,omitempty"`
	Description   string    `bson:"description,omitempty"`
	CreatedAt     time.Time `bson:"created_at,omitempty"`
	UpdatedAt     time.Time `bson:"updated_at,omitempty"`
}

func (i *InvoiceLog) CollectionName() string {
	return "invoiceLogs"
}
func (i *InvoiceLog) DatabaseName() string {
	return "db_superkul_order_test"
}

type InvoiceLogs interface {
	GetAll(Page int64, PageSize int64) ([]*InvoiceLog, int64, int64, error)
	FindById(Id string) (*InvoiceLog, error)
	FindByField(value string) []*InvoiceLog
	Delete(Id string) (*mongo.DeleteResult, error)
	Update(Id string, update map[string]interface{}) *InvoiceLog
}
