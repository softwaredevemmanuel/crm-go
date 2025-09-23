package admin

import (
	"fmt"
	"log"

	"github.com/xuri/excelize/v2"
	"crm-go/config"
	"crm-go/models"
	"github.com/gin-gonic/gin"
	
)
// ExportExcelHandler - download all data as Excel
func ExportExcelHandler(c *gin.Context) {
    filePath := "exported_data.xlsx"
    
    // Run the exporter
	ExportDataToExcel()

    // Serve file as download
    c.FileAttachment(filePath, "all-data.xlsx")
}

func ExportDataToExcel() {
	db := config.GetDB()

	// Example: Fetch categories
	var categories []models.Category
	if err := db.Find(&categories).Error; err != nil {
		log.Fatalf("failed to fetch categories: %v", err)
	}

	// Example: Fetch courses
	var courses []models.Course
	if err := db.Find(&courses).Error; err != nil {
		log.Fatalf("failed to fetch courses: %v", err)
	}

	// Create new Excel file
	f := excelize.NewFile()

	// ==============================
	// Sheet for Categories
	// ==============================
	catSheet := "Categories"
	index, _ := f.NewSheet(catSheet)
	f.SetCellValue(catSheet, "A1", "ID")
	f.SetCellValue(catSheet, "B1", "Name")
	f.SetCellValue(catSheet, "C1", "Created At")

	for i, cat := range categories {
		row := i + 2
		f.SetCellValue(catSheet, fmt.Sprintf("A%d", row), cat.ID)
		f.SetCellValue(catSheet, fmt.Sprintf("B%d", row), cat.Name)
		f.SetCellValue(catSheet, fmt.Sprintf("C%d", row), cat.CreatedAt)
	}

	// ==============================
	// Sheet for Courses
	// ==============================
	courseSheet := "Courses"
	f.NewSheet(courseSheet)
	f.SetCellValue(courseSheet, "A1", "ID")
	f.SetCellValue(courseSheet, "B1", "Title")
	f.SetCellValue(courseSheet, "C1", "CategoryID")
	f.SetCellValue(courseSheet, "D1", "Created At")

	for i, course := range courses {
		row := i + 2
		f.SetCellValue(courseSheet, fmt.Sprintf("A%d", row), course.ID)
		f.SetCellValue(courseSheet, fmt.Sprintf("B%d", row), course.Title)
		f.SetCellValue(courseSheet, fmt.Sprintf("C%d", row), course.Description)
		f.SetCellValue(courseSheet, fmt.Sprintf("D%d", row), course.CreatedAt)
	}

	// Set active sheet
	f.SetActiveSheet(index)

	// Save file
	if err := f.SaveAs("exported_data.xlsx"); err != nil {
		log.Fatalf("failed to save excel file: %v", err)
	}

	fmt.Println("Data exported successfully to exported_data.xlsx")
}
