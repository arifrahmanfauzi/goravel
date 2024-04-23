package services

import (
	"fmt"
	"github.com/goravel/framework/facades"
	"github.com/xuri/excelize/v2"
	"goravel/app/models"
)

type DownloadDropService struct {
}

func NewDownloadDropService() *DownloadDropService {
	return &DownloadDropService{}
}
func (Download *DownloadDropService) Download(Drops []*models.Drop) {
	xlsx := excelize.NewFile()
	var sheetName string = "Drops"
	err := xlsx.SetSheetName(xlsx.GetSheetName(0), sheetName)
	xlsx.SetCellValue(sheetName, "A1", "Trip Id")
	xlsx.SetCellValue(sheetName, "B1", "Job")
	xlsx.SetCellValue(sheetName, "C1", "Dispatch Number")
	for i, each := range Drops {
		xlsx.SetCellValue(sheetName, fmt.Sprintf("A%d", i+2), each.TripId)
		xlsx.SetCellValue(sheetName, fmt.Sprintf("B%d", i+2), each.Job)
		xlsx.SetCellValue(sheetName, fmt.Sprintf("C%d", i+2), each.DispatchNumber)
	}
	err = xlsx.SaveAs("./drops.xlsx")
	if err != nil {
		facades.Log().Error(err)
	}
}
