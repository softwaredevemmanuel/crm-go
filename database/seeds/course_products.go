package seeds

import (
	"log"

	"crm-go/config"
	"crm-go/models"

	"github.com/google/uuid"
)

func SeedCourseProductsTable() {
	db := config.GetDB()
	productID1, err1 := uuid.Parse("064d8891-7f46-41ff-af13-7f17f27ed11c")
	productID2, err2 := uuid.Parse("0b9e9362-9f1e-4261-bfc9-5975ba88067a")
	productID3, err3 := uuid.Parse("a903e5e5-02c3-453e-adbd-95b0c81fa974")
	productID4, err4 := uuid.Parse("a9f0edb1-d6db-4f8b-98ba-9cfc5563ecf1")
	productID5, err5 := uuid.Parse("d955fa42-1521-4af7-8687-55cf3a3be21d")

	courseID1, err6 := uuid.Parse("27d8ae14-4311-4380-8397-057ad5043fd6")
	courseID2, err7 := uuid.Parse("909b6026-30da-41f7-868f-42e6acba72c3")
	courseID3, err8 := uuid.Parse("b8ef3c14-d8ef-46fd-b63e-01b50cc9f227")
	courseID4, err9 := uuid.Parse("c40c9c00-0779-490a-931b-e8dbd91549bf")
	courseID5, err10 := uuid.Parse("cdbe9a63-c659-4912-abb0-58dcb9d2f341")
	
	courseProductID1, err11 := uuid.Parse("19a292e3-8376-4fd4-a1a1-04d6e21532b8")
	courseProductID2, err12 := uuid.Parse("417b2ded-bb70-4054-87e6-eb8e2abd3915")
	courseProductID3, err13 := uuid.Parse("4e84dfc2-aaca-4d9a-b96e-9074f294cae3")
	courseProductID4, err14 := uuid.Parse("78fcbd23-7414-4161-8a5a-dd89ec145b91")
	courseProductID5, err15 := uuid.Parse("834d0964-8552-4cb4-9f5f-dc8567ff63b7")
	courseProductID6, err16 := uuid.Parse("88ce907d-ae7f-462e-b224-c766e437f15a")
	courseProductID7, err17 := uuid.Parse("8b1a8f7e-3aa1-486c-99e3-00bdf6c37e19")
	courseProductID8, err18 := uuid.Parse("933e4d2b-f3f8-4358-9704-ca0b768d20f3")
	courseProductID9, err19 := uuid.Parse("ed9c891b-9bb3-4427-ae0d-d6292e1de5dc")
	courseProductID10, err20 := uuid.Parse("f7eb164d-9897-41d5-b233-4fb175144691")

if err1 != nil {
		log.Fatalf("❌ Invalid product UUID: %v", err1)
	}
	if err2 != nil {
		log.Fatalf("❌ Invalid product UUID: %v", err2)
	}
	if err3 != nil {
		log.Fatalf("❌ Invalid product UUID: %v", err3)
	}
	if err4 != nil {
		log.Fatalf("❌ Invalid product UUID: %v", err4)
	}
	if err5 != nil {
		log.Fatalf("❌ Invalid product UUID: %v", err5)
	}
	if err6 != nil {
		log.Fatalf("❌ Invalid course UUID: %v", err6)
	}
	if err7 != nil {
		log.Fatalf("❌ Invalid course UUID: %v", err7)
	}
	if err8 != nil {
		log.Fatalf("❌ Invalid course UUID: %v", err8)
	}
	if err9 != nil {
		log.Fatalf("❌ Invalid course UUID: %v", err9)
	}
	if err10 != nil {
		log.Fatalf("❌ Invalid course UUID: %v", err10)
	}
if err11 != nil {
		log.Fatalf("❌ Invalid course product UUID: %v", err11)
	}
	if err12 != nil {
		log.Fatalf("❌ Invalid course product UUID: %v", err12)
	}
	if err13 != nil {
		log.Fatalf("❌ Invalid course product UUID: %v", err13)
	}
	if err14 != nil {
		log.Fatalf("❌ Invalid course product UUID: %v", err14)
	}
	if err15 != nil {
		log.Fatalf("❌ Invalid course product UUID: %v", err15)
	}
	if err16 != nil {
		log.Fatalf("❌ Invalid course product UUID: %v", err16)
	}
	if err17 != nil {
		log.Fatalf("❌ Invalid course product UUID: %v", err17)
	}
	if err18 != nil {
		log.Fatalf("❌ Invalid course product UUID: %v", err18)
	}
	if err19 != nil {
		log.Fatalf("❌ Invalid course product UUID: %v", err19)
	}
	if err20 != nil {
		log.Fatalf("❌ Invalid course product UUID: %v", err20)
	}



	course_products := []models.CourseProductTable{
		{
		ID: 	  courseProductID1,
		CourseID: courseID1,
		ProductID: productID3,
		},
		{
		ID: 	  courseProductID2,
		CourseID: courseID2,
		ProductID: productID4,	
		},
		{
		ID: 	  courseProductID3,
		CourseID: courseID3,
		ProductID: productID1,	
		},
		{
		ID: 	  courseProductID4,
		CourseID: courseID4,
		ProductID: productID2,	
		},
		{
		ID: 	  courseProductID5,
		CourseID: courseID5,
		ProductID: productID5,	
		},	
		{
		ID: 	  courseProductID6,
		CourseID: courseID1,
		ProductID: productID1,
		},	
		{
		ID: 	  courseProductID7,
		CourseID: courseID2,
		ProductID: productID2,
		},
		{
		ID: 	  courseProductID8,
		CourseID: courseID3,
		ProductID: productID3,
		},	
		{
		ID: 	  courseProductID9,
		CourseID: courseID4,
		ProductID: productID4,
		},
		{
		ID: 	  courseProductID10,
		CourseID: courseID5,
		ProductID: productID5,
		},	
	
	

		
		
	}

	

	for _, course_product := range course_products {
		if err := db.Create(&course_product).Error; err != nil {
			log.Printf("❌ Failed to seed Course Product Table: %v", err)
		} else {
			log.Printf("✅ Course Product Seeded Successfully:")
		}
	}
}

