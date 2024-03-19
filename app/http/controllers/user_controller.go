package controllers

import (
	"github.com/goravel/framework/contracts/http"
)

type UserController struct {
	//Dependent services
}

func NewUserController() *UserController {
	return &UserController{
		//Inject services
	}
}

func (r *UserController) Show(ctx http.Context) http.Response {
	url := ctx.Request().Url()
	id := ctx.Request().RouteInt("id")
	return ctx.Response().Success().Json(http.Json{
		"Hello": "Goravel",
		"url":   url,
		"id":    id,
	})
}
func (r *UserController) Store(ctx http.Context) http.Response {
	// name := ctx.Request().Input("name")
	var user User
	err := ctx.Request().Bind(&user)
	if err != nil {
		return ctx.Response().Status(http.StatusBadRequest).Json(http.Json{
			"status": "Error",
		})
	}
	return ctx.Response().Success().Json(http.Json{
		"name": user.Name,
		"type": user.Type,
	})
}

type User struct {
	Name  string `form:"name" json:"name"`
	Email string `form:"email" json:"email"`
	Type  string `form:"type" json:"type"`
}
