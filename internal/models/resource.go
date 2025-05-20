package models

import (
	"time"

	"github.com/google/uuid"
)

type ClientOrganizationUnit struct {
	ResourceID       uuid.UUID `json:"resource_id"`
	ParentResourceID uuid.UUID `json:"parent_resource_id"`
	ResourceTypeID   uuid.UUID `json:"resource_type_id,omitempty"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	RowStatus        int       `json:"row_status"`
	CreatedBy        uuid.UUID `json:"created_by"`
	UpdatedBy        uuid.UUID `json:"updated_by"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at,omitempty"`
}
