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
