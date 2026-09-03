package seeds

import (
	"log"

	"crm-go/config"
	"crm-go/models"
	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

func SeedDepartments() error {
	db := config.GetDB()

	adminID := uuid.MustParse("fe4547a7-4c81-4bc2-bc81-5bbbce2fb5bd")
	headID := uuid.MustParse("5a853260-31fc-44ee-9d69-bb2a2957ba48")

	LanguagesDepartment := uuid.MustParse("4a547934-2c83-4705-8616-e1e817b66bb2")
	MathematicsDepartment := uuid.MustParse("78a57227-d318-4264-b373-c0ae96f7e83b")
	SciencesDepartment := uuid.MustParse("842dfd3f-6344-4c80-89a9-72ce9bb3ab15")
	SocialSciencesDepartment := uuid.MustParse("8bc3ebc8-79b7-4092-931f-d235e3bd7d1f")
	ArtsDepartment := uuid.MustParse("94d593bf-3ef4-4c10-9a82-ac32603a4350")
	TechnologyDepartment := uuid.MustParse("a9dddddb-28b8-4264-8bef-675fe5cc4e5c")
	BusinessStudiesDepartment := uuid.MustParse("b323135b-b1bc-4d71-9fa0-3d6affd8ad8b")
	AgriculturalSciencesDepartment := uuid.MustParse("d6ec9227-b333-4c58-9829-060005320746")
	HomeEconomicsDepartment := uuid.MustParse("eaa85465-26f5-430d-a4d8-71150b1352e7")
	ReligiousStudiesDepartment := uuid.MustParse("f7329d7c-1fce-4737-b117-7976729877de")
	PhysicalEducationDepartment := uuid.MustParse("ca1d4208-32d5-4dbf-9a5d-fbc8089dc341")

	departments := []models.Department{
		{
			ID:          LanguagesDepartment,
			Name:        "Languages",
			Code:        "LANG",
			Description: "Covers English, French, Literature and other language subjects.",
			HeadOfDept:  &headID,
			Status:      "active",
			CreatedBy:   adminID,
		},
		{
			ID:          MathematicsDepartment,
			Name:        "Mathematics",
			Code:        "MATH",
			Description: "Covers Mathematics and advanced mathematical studies.",
			HeadOfDept:  &headID,
			Status:      "active",
			CreatedBy:   adminID,
		},
		{
			ID:          SciencesDepartment,
			Name:        "Sciences",
			Code:        "SCI",
			Description: "Covers Biology, Chemistry, Physics and Basic Science.",
			HeadOfDept:  &headID,
			Status:      "active",
			CreatedBy:   adminID,
		},
		{
			ID:          SocialSciencesDepartment,
			Name:        "Social Sciences",
			Code:        "SOC",
			Description: "Covers Economics, Geography, Government and related subjects.",
			HeadOfDept:  &headID,
			Status:      "active",
			CreatedBy:   adminID,
		},
		{
			ID:          ArtsDepartment,
			Name:        "Arts",
			Code:        "ART",
			Description: "Covers Visual Arts, Music, Drama and Creative Arts.",
			HeadOfDept:  &headID,
			Status:      "active",
			CreatedBy:   adminID,
		},
		{
			ID:          TechnologyDepartment,
			Name:        "Technology",
			Code:        "TECH",
			Description: "Covers Computer Studies, Basic Technology and Technical Drawing.",
			HeadOfDept:  &headID,
			Status:      "active",
			CreatedBy:   adminID,
		},
		{
			ID:          BusinessStudiesDepartment,
			Name:        "Business Studies",
			Code:        "BUS",
			Description: "Covers Accounting, Commerce, Marketing and business-related studies.",
			HeadOfDept:  &headID,
			Status:      "active",
			CreatedBy:   adminID,
		},
		{
			ID:          AgriculturalSciencesDepartment,
			Name:        "Agricultural Sciences",
			Code:        "AGR",
			Description: "Covers Agricultural Science, farming and related agricultural studies.",
			HeadOfDept:  &headID,
			Status:      "active",
			CreatedBy:   adminID,
		},
		{
			ID:          HomeEconomicsDepartment,
			Name:        "Home Economics",
			Code:        "HEC",
			Description: "Covers Food & Nutrition, clothing, family living and home management.",
			HeadOfDept:  &headID,
			Status:      "active",
			CreatedBy:   adminID,
		},
		{
			ID:          ReligiousStudiesDepartment,
			Name:        "Religious Studies",
			Code:        "REL",
			Description: "Covers Christian Religious Studies and Islamic Religious Studies.",
			HeadOfDept:  &headID,
			Status:      "active",
			CreatedBy:   adminID,
		},
		{
			ID:          PhysicalEducationDepartment,
			Name:        "Physical Education",
			Code:        "PHE",
			Description: "Covers physical fitness, sports, health education, recreation and personal wellbeing.",
			HeadOfDept:  &headID,
			Status:      "active",
			CreatedBy:   adminID,
		},
	}

	for _, department := range departments {
		result := db.Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "id"},
			},
			DoUpdates: clause.AssignmentColumns([]string{
				"name",
				"code",
				"description",
				"head_of_dept",
				"status",
				"created_by",
			}),
		}).Create(&department)

		if result.Error != nil {
			log.Printf(
				"❌ Failed to seed department %s: %v",
				department.Name,
				result.Error,
			)
		} else {
			log.Printf(
				"✅ Seeded/updated department: %s (%s)",
				department.Name,
				department.Code,
			)
		}
	}
	return nil
}
