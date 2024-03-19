package routes

import (
	"github.com/goravel/framework/facades"

	"goravel/app/http/controllers"
	"goravel/app/http/middleware"
)

func Api() {
	userController := controllers.NewUserController()
	blogController := controllers.NewBlogController()
	verifyController := controllers.NewVerifyController()
	tripController := controllers.NewTripController()
	orderController := controllers.NewOrderController()
	//dropController := controllers.NewDropController()
	facades.Route().Get("/users/{id}", userController.Show)
	facades.Route().Post("/users", userController.Store)
	facades.Route().Middleware(middleware.Verify()).Get("/blogs", blogController.Index)
	facades.Route().Get("/blogs/{id}", blogController.Find)
	facades.Route().Get("/verify", verifyController.Index)
	facades.Route().Get("/trips", tripController.Index)
	facades.Route().Get("/orders", orderController.GetAll)
	facades.Route().Get("/drops", controllers.NewDropController().Index)
	facades.Route().Patch("/drops/{id}", controllers.NewDropController().Update)
	facades.Route().Delete("/drops/{id}", controllers.NewDropController().Delete)
	facades.Route().Post("/invoice-log", controllers.NewInvoiceLogController().Create)
	facades.Route().Patch("/invoice-log/{id}", controllers.NewInvoiceLogController().Update)
	facades.Route().Get("/invoice-log/{id}", controllers.NewInvoiceLogController().FindById)
	facades.Route().Delete("/invoice-log/{id}", controllers.NewInvoiceLogController().Delete)
	facades.Route().Get("/invoice-log/invoiceNumber", controllers.NewInvoiceLogController().FindByField)
	facades.Route().Get("/invoice-log", controllers.NewInvoiceLogController().Index)
}
