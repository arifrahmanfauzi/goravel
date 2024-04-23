package models

type CorporateInvoice struct {
	ID            string `bson:"_id,omitempty"`
	CustomerType  string `bson:"customerType,omitempty"`
	CustomerName  string `bson:"customerName,omitempty"`
	InvoiceType   string `bson:"invoiceType,omitempty"`
	InvoiceStatus string `bson:"invoiceStatus,omitempty"`
}

func (receiver CorporateInvoice) DatabaseName() string {
	return "db_superkul_order"
}
func (receiver CorporateInvoice) CollectionName() string {
	return "corporateInvoice"
}
