package models

type CustomerOrder struct {
	ID                string              `bson:"_id,omitempty"`
	ServiceName       string              `bson:"serviceName,omitempty"`
	OrderNumber       string              `bson:"orderNumber,omitempty"`
	InvoiceNumber     string              `bson:"invoiceNumber,omitempty"`
	AdditionalService []AdditionalService `bson:"additionalService,omitempty"`
}

type AdditionalService struct {
	Name        string `bson:"name"`
	Price       int    `bson:"price"`
	IsMandatory int    `bson:"is_mandatory"`
}

func (co CustomerOrder) CollectionName() string {
	return "customerOrder"
}
func (co CustomerOrder) DatabaseName() string {
	return "db_superkul_order"
}
