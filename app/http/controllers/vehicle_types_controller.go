package controllers

import (
	"github.com/goravel/framework/contracts/http"
)

type VehicleTypesController struct {
	//Dependent services
}

func NewVehicleTypesController() *VehicleTypesController {
	return &VehicleTypesController{
		//Inject services
	}
}

func (r *VehicleTypesController) Index(ctx http.Context) http.Response {
	return nil
}
