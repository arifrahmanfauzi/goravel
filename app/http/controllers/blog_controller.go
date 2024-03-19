package controllers

import (
	"goravel/app/models"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

type BlogController struct {
	//Dependent services
}

func NewBlogController() *BlogController {
	return &BlogController{
		//Inject services
	}
}

func (r *BlogController) Index(ctx http.Context) http.Response {
	var blogs []models.Blog
	var total int64
	if err := facades.Orm().Connection("mysql").Query().Paginate(1, 15, &blogs, &total); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"status":  "Error",
			"message": err.Error(),
		})

	}
	return ctx.Response().Success().Json(http.Json{
		"status": "OK",
		"data":   blogs,
		"total":  total,
	})
}
func (r *BlogController) Find(ctx http.Context) http.Response {
	var blogs models.Blog
	if err := facades.Orm().Query().Find(&blogs, ctx.Request().RouteInt64("id")); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"status":  "Error",
			"message": err.Error(),
		})
	}
	return ctx.Response().Success().Json(http.Json{
		"status": "OK",
		"data":   blogs,
	})
}
