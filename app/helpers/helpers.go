package helpers

func paginate(Page int, TotalRecord int, PageSize int) (int, int) {
	totalPages := TotalRecord / PageSize
	if TotalRecord%PageSize != 0 {
		totalPages++
	}
	skip := (Page - 1) * PageSize
	return skip, totalPages
}
