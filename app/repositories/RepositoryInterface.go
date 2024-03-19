package repositories

type Repository interface {
	GetAll(page int64, pageSize int64, total *int64, totalPage *int64)
	FindById(ID string)
	Create()
	Update()
	Delete()
}
