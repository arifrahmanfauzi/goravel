package controllers

import (
	"github.com/goravel/framework/contracts/http"
)

type VerifyController struct {
	//Dependent services
}

func NewVerifyController() *VerifyController {
	return &VerifyController{
		//Inject services
	}
}

func (r *VerifyController) Index(ctx http.Context) http.Response {

	token := ctx.Request().Header("Authorization", "")
	return ctx.Response().Success().Json(http.Json{
		"token": token,
	})
}
