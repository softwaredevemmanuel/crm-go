package seeds

import (
	"log"
	"time"

	"crm-go/config"
	"crm-go/models"

	"github.com/google/uuid"
)

func SeedProducts() {
	db := config.GetDB()
	productID1, err1 := uuid.Parse("064d8891-7f46-41ff-af13-7f17f27ed11c")
	productID2, err2 := uuid.Parse("0b9e9362-9f1e-4261-bfc9-5975ba88067a")
	productID3, err3 := uuid.Parse("a903e5e5-02c3-453e-adbd-95b0c81fa974")
	productID4, err4 := uuid.Parse("a9f0edb1-d6db-4f8b-98ba-9cfc5563ecf1")
	productID5, err5 := uuid.Parse("d955fa42-1521-4af7-8687-55cf3a3be21d")

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

	products := []models.Product{
		{
			ID:          productID1,
			Name:       "HP Laptop 15s-fq2713TU",
			Description: "High performance laptop with Intel i5, 8GB RAM, 512GB SSD.",
			Price: 300000.00,
			CompareAtPrice: 350000.00,
			Image:       "https://example.com/images/laptop.png",
			RequiresShipping:    false,
			Status:     "active",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          productID2,
			Name:       "Apple iPhone 13",
			Description: "Latest Apple iPhone with A15 Bionic chip, 128GB storage.",
			Price: 800000.00,
			CompareAtPrice: 900000.00,
			Image:       "https://example.com/images/iphone13.png",
			RequiresShipping:    false,
			Status:     "active",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          productID3,
			Name:       "Samsung Galaxy S21",
			Description: "Flagship Samsung phone with excellent camera and display.",
			Price: 700000.00,
			CompareAtPrice: 750000.00,
			Image:       "https://example.com/images/galaxys21.png",
			RequiresShipping:    false,
			Status:     "active",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          productID4,
			Name:       "Sony WH-1000XM4 Headphones",
			Description: "Industry leading noise cancelling over-ear headphones.",
			Price: 250000.00,
			CompareAtPrice: 300000.00,
			Image:       "https://example.com/images/sonyheadphones.png",
			RequiresShipping:    false,
			Status:     "active",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          productID5,
			Name:       "Dell UltraSharp U2723QE Monitor",
			Description: "27-inch 4K UHD monitor with excellent color accuracy.",
			Price: 400000.00,
			CompareAtPrice: 450000.00,
			Image:       "https://example.com/images/dellmonitor.png",
			RequiresShipping:    false,
			Status:     "active",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),	
		},
	}

	for _, product := range products {
		if err := db.Create(&product).Error; err != nil {
			log.Printf("❌ Failed to seed Product: %v", err)
		} else {
			log.Printf("✅ Seeded product: %v", &product.Name)
		}
	}
}

