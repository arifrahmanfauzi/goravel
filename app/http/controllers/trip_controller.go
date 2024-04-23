package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"goravel/app/repositories"
	"goravel/app/transformers"
)

type TripController struct {
	//Dependent services
	Repositories *repositories.TripRepository
}

func NewTripController() *TripController {
	return &TripController{
		//Inject services
		Repositories: repositories.NewTripRepository(),
	}
}

func (r *TripController) Index(ctx http.Context) http.Response {
	page := ctx.Request().QueryInt64("page", 1)
	var total int64
	var totalPage int64
	trips, err := r.Repositories.GetAll(page, 15, &total, &totalPage)
	if err != nil {
		facades.Log().Error(err)
	}
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
func (r *TripController) Find(ctx http.Context) http.Response {
	trips := r.Repositories.FindTripByOrderNumber(ctx.Request().Route("id"))
	return ctx.Response().Json(http.StatusOK, http.Json{
		"data": trips,
	})
}
