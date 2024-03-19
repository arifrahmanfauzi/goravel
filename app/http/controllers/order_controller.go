package controllers

import (
	"context"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"goravel/app/repositories"
	"goravel/database/mongodb"
)

type OrderController struct {
	//Dependent services
}

func NewOrderController() *OrderController {
	return &OrderController{
		//Inject services
	}
}

func (r *OrderController) Index(ctx http.Context) http.Response {
	return nil
}
func (r OrderController) GetAll(ctx http.Context) http.Response {
	//clientOptions := options.Client().ApplyURI(facades.Config().GetString("DB_STRING", ""))
	// Connect to MongoDB
	//client, err := mongo.Connect(context.Background(), clientOptions)
	client, err := mongodb.Init()
	orderRepository := repositories.NewOrderRepository(client)
	if err != nil {
		facades.Log().Error(err)
	}
	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		facades.Log().Error(err)
	}
	page := ctx.Request().QueryInt64("page", 1)
	var total int64
	var totalPage int64
	orders, err := orderRepository.GetAll(page, 15, &total, &totalPage)
	return ctx.Response().Json(http.StatusOK, http.Json{
		"Status": "OK",
		"Data":   orders,
	})
}
