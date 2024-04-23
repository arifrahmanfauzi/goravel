package controllers

import (
	"fmt"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"goravel/app/models"
	"goravel/app/repositories"
	"goravel/app/transformers"
)

type InvoiceLog struct {
	//Dependent services
	IR *repositories.InvoiceLogRepository
}

func NewInvoiceLogController() *InvoiceLog {
	return &InvoiceLog{
		//Inject services
		IR: repositories.NewInvoiceLogRepository(),
	}
}

func (r *InvoiceLog) Create(ctx http.Context) http.Response {
	var InvoiceLogs models.InvoiceLog
	err := ctx.Request().Bind(&InvoiceLogs)
	if err != nil {
		facades.Log().Error(err)
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"message": err,
		})
	}
	result, err := r.IR.Create(&InvoiceLogs)
	if err != nil {
		facades.Log().Error(err)
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"message": err,
		})
	}
	return ctx.Response().Json(http.StatusOK, http.Json{
		"status": "OK",
		"data":   InvoiceLogs,
		"raw":    result,
	})
}
func (r *InvoiceLog) Index(ctx http.Context) http.Response {
	page := ctx.Request().QueryInt64("page", 1)
	invoiceLogs, totalRecord, totalPage, err := r.IR.GetAll(page, 15)
	if err != nil {
		facades.Log().Error(err)
	}
	return ctx.Response().Json(http.StatusAccepted, transformers.Response{
		Status: "OK",
		Data:   invoiceLogs,
		Meta: transformers.Meta{
			Pagination: transformers.Pagination{
				Total:       totalRecord,
				Count:       15,
				PerPage:     15,
				CurrentPage: page,
				TotalPages:  totalPage,
			},
		},
	})
}
func (r *InvoiceLog) FindById(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")
	//invoiceLogRepository := repositories.NewInvoiceLogRepository()
	result, err := r.IR.FindById(id)
	if err != nil {
		facades.Log().Error(err)
		return nil
	}
	return ctx.Response().Json(http.StatusAccepted, transformers.Response{
		Status: "OK",
		Data:   result,
		Meta: transformers.Meta{
			Pagination: transformers.Pagination{},
		},
	})
}
func (r *InvoiceLog) FindByField(ctx http.Context) http.Response {
	invoiceLogRepository := repositories.NewInvoiceLogRepository()
	invoiceLogs := invoiceLogRepository.FindByField("INV-SD-20230525353352")
	return ctx.Response().Json(http.StatusOK, http.Json{
		"data": invoiceLogs,
	})
}
func (r *InvoiceLog) Update(ctx http.Context) http.Response {
	//invoiceLogRepository := repositories.NewInvoiceLogRepository()

	var InvoiceLogs map[string]any
	err := ctx.Request().Bind(&InvoiceLogs)
	if err != nil {
		facades.Log().Error(err)
	}
	invoiceLog := r.IR.Update(ctx.Request().Route("id"), InvoiceLogs)
	return ctx.Response().Json(http.StatusOK, http.Json{
		"status": "OK",
		"data":   InvoiceLogs,
		"raw":    invoiceLog,
	})
}
func (r *InvoiceLog) Delete(ctx http.Context) http.Response {
	//invoiceLogRepository := repositories.NewInvoiceLogRepository()
	deleteResult, err := r.IR.Delete(ctx.Request().Route("id"))
	if err != nil {
		facades.Log().Error(err)
		return nil
	}
	return ctx.Response().Json(http.StatusOK, http.Json{
		"status":  "OK",
		"message": fmt.Sprintf("rows affected %d", deleteResult.DeletedCount),
	})
}
