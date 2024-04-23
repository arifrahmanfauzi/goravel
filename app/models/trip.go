package models

type Trip struct {
	ID         string `bson:"_id,omitempty"`
	OrderId    string `bson:"orderId,omitempty"`
	TripNumber string `bson:"tripNumber,omitempty"`
	Distance   int64  `bson:"distance,omitempty"`
	Pick       []any  `bson:"pick,omitempty"`
}

func (tr *Trip) CollectionName() string {
	return "customerTripPlanning"
}
func (tr *Trip) DatabaseName() string {
	return "db_superkul_order"
}

type TripInterface interface {
	Create(trip *Trip) error
	FindByID(id string) (*Trip, error)
	GetAll(page, pageSize int64, total *int64, totalPage *int64) ([]*Trip, error)
	FindTripByOrderNumber(OrderNumber string) []*Trip
}
