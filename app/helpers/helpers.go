package helpers

// Paginate is helper to produce pagination
// Page int is current page
//
//	TotalRecord int is total a whole document which has been query
//
// PageSize int is Limit number of document to display
func Paginate(Page int, TotalRecord int, PageSize int) (int, int) {
	//Calculate total Pages from count record
	TotalPages := TotalRecord / PageSize
	if TotalRecord%PageSize != 0 {
		TotalPages++
	}
	//Calculate Number of Document to skip
	NumberOfDocumentSkip := (Page - 1) * PageSize
	return NumberOfDocumentSkip, TotalPages
}
