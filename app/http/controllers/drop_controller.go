package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"goravel/app/repositories"
	"goravel/app/transformers"
)

type DropController struct {
	//Dependent services
	Repositories *repositories.DropRepository
}

func NewDropController() *DropController {
	return &DropController{
		//Inject services
		Repositories: repositories.NewDropRepository(),
	}
}

func (r *DropController) Index(ctx http.Context) http.Response {
	page := ctx.Request().QueryInt64("page", 1)
	dropRepository := repositories.NewDropRepository()
	var total int64
	var totalPage int64
	drops, err := dropRepository.GetAll(page, 15, &total, &totalPage)
	if err != nil {
		facades.Log().Error(err)
	}
	return ctx.Response().Json(http.StatusAccepted, transformers.Response{
		Status: "OK",
		Data:   drops,
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

func (r *DropController) Update(ctx http.Context) http.Response {
	dropRepository := repositories.NewDropRepository()
	id := ctx.Request().Route("id")

	var drop map[string]any
	err := ctx.Request().Bind(&drop)
	if err != nil {
		facades.Log().Error(err)
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"status": "Failed",
			"data":   err,
		})
	}
	dropRepository.Update(id, drop)
	return ctx.Response().Json(http.StatusOK, http.Json{
		"status": "OK",
		"data":   drop,
	})
}

func (r *DropController) Delete(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")
	res, err := repositories.NewDropRepository().Delete(id)
	if err != nil {
		facades.Log().Error(err)
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"status":  "Error",
			"message": "Failed to delete record",
		})
	}
	if res.DeletedCount < 0 {
		return ctx.Response().Json(http.StatusOK, http.Json{
			"status":  "OK",
			"message": "Success Deleted!",
		})
	}
	return ctx.Response().Json(http.StatusOK, http.Json{
		"status":  "OK",
		"message": "No record has been deleted!",
	})

}

func (r *DropController) Find(ctx http.Context) http.Response {
	page := ctx.Request().QueryInt64("page", 1)
	res, total, totalPage, err := r.Repositories.FindByTripNumber(ctx.Request().Route("id"), page, 15)
	if err != nil {
		facades.Log().Error(err)
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"status":  "ERROR",
			"message": err,
		})
	}
	return ctx.Response().Json(http.StatusOK, transformers.Response{
		Status: "OK",
		Data:   res,
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
