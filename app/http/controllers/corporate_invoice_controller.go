package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"goravel/app/repositories"
)

type CorporateInvoiceController struct {
	//Dependent services
	Repository *repositories.CorporateInvoiceRepository
}

func NewCorporateInvoiceController() *CorporateInvoiceController {
	return &CorporateInvoiceController{
		//Inject services
		Repository: repositories.NewCorporateInvoiceRepository(),
	}
}

func (r *CorporateInvoiceController) Index(ctx http.Context) http.Response {
	return nil
}
func (r *CorporateInvoiceController) GetAll(ctx http.Context) http.Response {
	var Pages = ctx.Request().QueryInt64("page", 1)
	var Limit = ctx.Request().QueryInt64("limit", int64(facades.Config().GetInt("app.pagination", 15)))
	var InvoiceNumber = ctx.Request().Query("invoiceNumber", "")

	corporateInvoice, Pagination := r.Repository.GetAll(Pages, Limit, InvoiceNumber)
	return ctx.Response().Json(http.StatusOK, http.Json{
		"data": corporateInvoice,
		"meta": map[string]any{
			"pagination": Pagination,
		},
	})
}
