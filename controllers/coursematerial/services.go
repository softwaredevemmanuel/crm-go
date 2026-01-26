package coursematerial

import (
	"crm-go/models"
	"encoding/json"
	"gorm.io/gorm"
	"github.com/google/uuid"
)

// DeleteWithArchive moved to services/delete_service.go to avoid duplicate declarations
// (implementation exists later in this file under the services/delete_service.go section)

// services/delete_service.go
func DeleteWithArchive(
	tx *gorm.DB,
	entityType string,
	entityID uuid.UUID,
	data any, // This is for archiving only, NOT for deletion
	deletedBy *uuid.UUID,
	reason string,
	model interface{}, // Add this parameter - the actual model to delete
) error {

	// First, marshal the data for archiving
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Create the archive record
	archive := models.DeletedRecord{
		EntityType: entityType,
		EntityID:   entityID,
		Data:       jsonData,
		DeletedBy:  deletedBy,
		Reason:     reason,
	}

	if err := tx.Create(&archive).Error; err != nil {
		return err
	}

	// Delete the actual model, NOT the data parameter
	if err := tx.Delete(model, "id = ?", entityID).Error; err != nil {
		return err
	}

	return nil
}