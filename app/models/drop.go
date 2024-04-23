package models

import "go.mongodb.org/mongo-driver/mongo"

type Drop struct {
	ID             string `bson:"_id,omitempty"`
	TripId         string `bson:"tripId,omitempty"`
	Job            string `bson:"job,omitempty"`
	DispatchNumber string `bson:"dispatchNumber,omitempty"`
}

<<<<<<< HEAD
func (Drop Drop) CollectionName() string {
	return "customerTripPlanningDt"
}
func (Drop Drop) DatabaseName() string {
	return "db_superkul_order"
=======
func (dop *Drop) CollectionName() string {
	return "customerTripPlanningDt"
}
func (dop *Drop) DatabaseName() string {
	return "db_superkul_order_test"
>>>>>>> origin/main
}

type DropInterface interface {
	GetAll(page int64, pageSize int64, total *int64, totalPage *int64) ([]*Drop, error)
	Update(ID string, update map[string]interface{}) *Drop
	Delete(ID string) (*mongo.DeleteResult, error)
	FindByTripNumber(TripNumber string, Page int64, PageSize int64) ([]*Drop, int64, int64, error)
}
