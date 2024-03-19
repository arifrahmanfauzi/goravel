package controllers

import (
	"context"
	"fmt"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"goravel/app/repositories"
	"goravel/app/transformers"
)

type TripController struct {
	//Dependent services
}

func NewTripController() *TripController {
	return &TripController{
		//Inject services
	}
}

func (r *TripController) Index(ctx http.Context) http.Response {
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
	fmt.Println("Connected to MongoDB!")
	page := ctx.Request().QueryInt64("page", 1)
	tripRepository := repositories.NewTripRepository(client, "db_superkul_order_test", "customerTripPlanning")
	var total int64
	var totalPage int64
	trips, err := tripRepository.GetAll(page, 15, &total, &totalPage)

	return ctx.Response().Json(http.StatusAccepted, transformers.Response{
		Status: "OK",
		Data:   trips,
		Meta: transformers.Meta{
			Pagination: transformers.Pagination{
				Total:       total,
				Count:       15,
				PerPage:     15,
				CurrentPage: page,
				TotalPages:  totalPage,
			},
		},
	})
}
