package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TenantResource struct {
	ResourceID       uuid.UUID  `gorm:"type:char(36);primaryKey;column:resource_id" json:"resource_id"`
	ParentResourceID *uuid.UUID `gorm:"type:char(36);column:parent_resource_id" json:"parent_resource_id"`
	ResourceTypeID   uuid.UUID  `gorm:"type:char(36);not null;column:resource_type_id" json:"resource_type_id"` // foreign key to resource_type
	Name             string     `gorm:"size:45;not null;column:name" json:"name"`
	TenantID         *uuid.UUID `gorm:"type:char(36);column:tenant_id" json:"tenant_id"`
	RowStatus        int        `gorm:"default:1;column:row_status" json:"row_status"`
	CreatedBy        uuid.UUID  `gorm:"size:45;column:created_by" json:"created_by"`
	UpdatedBy        uuid.UUID  `gorm:"size:45;column:updated_by" json:"updated_by"`
	CreatedAt        time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"column:updated_at;autoCreateTime" json:"updated_at"`
}

func (t *TenantResource) TableName() string {
	return "tnt_resources"
}

type Mst_ResourceTypes struct {
	ResourceTypeID uuid.UUID `gorm:"type:char(36);primaryKey;column:resource_type_id" json:"resource_type_id"`
	ServiceID      uuid.UUID `gorm:"type:char(36);not null;column:service_id" json:"service_id"`
	Name           string    `gorm:"size:45;not null;column:name" json:"name"`
	RowStatus      int       `gorm:"default:1;column:row_status" json:"row_status"`
	CreatedBy      uuid.UUID `gorm:"size:45;column:created_by" json:"created_by"`
	UpdatedBy      uuid.UUID `gorm:"size:45;column:updated_by" json:"updated_by"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoCreateTime" json:"updated_at"`
}

func (t *Mst_ResourceTypes) TableName() string {
	return "mst_resource_types"
}

// Tenant struct aligned with schema
type TenantMetadata struct {
	ResourceID uuid.UUID       `gorm:"type:char(36);not null" json:"resource_id"`
	Metadata   json.RawMessage `gorm:"type:json;" json:"metadata"`
	RowStatus  int             `gorm:"default:1;column:row_status" json:"row_status"`
	CreatedBy  uuid.UUID       `gorm:"size:45" json:"created_by"`
	UpdatedBy  uuid.UUID       `gorm:"size:45" json:"updated_by"`
	CreatedAt  time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
}

// BeforeCreate hook to generate UUID before saving
func (t *TenantMetadata) BeforeCreate(tx *gorm.DB) (err error) {
	//t.ID = uuid.New()
	return
}

func (t *TenantMetadata) TableName() string {
	return "tnt_resources_metadata"
}

type TenantRoleAssignments struct {
	ResourceID  uuid.UUID `gorm:"type:char(36);primaryKey;column:resource_id" json:"resource_id"`
	Name        string    `gorm:"size:45;not null;column:name" json:"name"`
	Version     string    `gorm:"size:45;not null;column:version" json:"version"`
	PrincipalID uuid.UUID `gorm:"type:char(36);column:principal_id" json:"principal_id"`
	RoleID      uuid.UUID `gorm:"type:char(36);column:role_id" json:"role_id"`
	RowStatus   int       `gorm:"default:1;column:row_status" json:"row_status"`
	CreatedBy   uuid.UUID `gorm:"size:45;column:created_by" json:"created_by"`
	UpdatedBy   uuid.UUID `gorm:"size:45;column:updated_by" json:"updated_by"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoCreateTime" json:"updated_at"`
}

func (t *TenantRoleAssignments) TableName() string {
	return "tnt_role_assignments"
}

type TenantPrincipals struct {
	ResourceID      uuid.UUID       `gorm:"type:char(36);primaryKey;column:resource_id" json:"resource_id"`
	PrincipalTypeID uuid.UUID       `gorm:"size:45;not null;column:principal_type_id" json:"principal_type_id"`
	Name            string          `gorm:"size:45;not null;column:name" json:"name"`
	Email           string          `gorm:"size:45;not null;column:email" json:"email"`
	Metadata        json.RawMessage `gorm:"type:json;" json:"metadata"`
	RowStatus       int             `gorm:"default:1;column:row_status" json:"row_status"`
	CreatedBy       uuid.UUID       `gorm:"size:45;column:created_by" json:"created_by"`
	UpdatedBy       uuid.UUID       `gorm:"size:45;column:updated_by" json:"updated_by"`
	CreatedAt       time.Time       `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time       `gorm:"column:updated_at;autoCreateTime" json:"updated_at"`
}

func (t *TenantPrincipals) TableName() string {
	return "tnt_principals"
}

type TenantRoles struct {
	ResourceID     uuid.UUID `gorm:"type:char(36);primaryKey;column:resource_id" json:"resource_id"`
	ResourceTypeID uuid.UUID `gorm:"size:45;not null;column:resource_type_id" json:"resource_type_id"`
	RoleType       string    `gorm:"size:45;not null;column:role_type" json:"role_type"`
	Name           string    `gorm:"size:45;not null;column:name" json:"name"`
	Version        string    `gorm:"size:45;not null;column:version" json:"version"`
	RowStatus      int       `gorm:"default:1;column:row_status" json:"row_status"`
	CreatedBy      uuid.UUID `gorm:"size:45;column:created_by" json:"created_by"`
	UpdatedBy      uuid.UUID `gorm:"size:45;column:updated_by" json:"updated_by"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoCreateTime" json:"updated_at"`
}

func (t *TenantRoles) TableName() string {
	return "tnt_roles"
}

type TenantRolePermissions struct {
	RolePermissionID uuid.UUID `gorm:"type:char(36);primaryKey;column:role_permission_id" json:"role_permission_id"`
	RoleID           uuid.UUID `gorm:"size:45;not null;column:role_id" json:"role_id"`
	PermissionID     uuid.UUID `gorm:"size:45;not null;column:permission_id" json:"permission_id"`
	RowStatus        int       `gorm:"default:1;column:row_status" json:"row_status"`
	CreatedBy        uuid.UUID `gorm:"size:45;column:created_by" json:"created_by"`
	UpdatedBy        uuid.UUID `gorm:"size:45;column:updated_by" json:"updated_by"`
	CreatedAt        time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at;autoCreateTime" json:"updated_at"`
}

func (t *TenantRolePermissions) TableName() string {
	return "tnt_role_permissions"
}

type MstPrincipalTypes struct {
	PrincipalTypeID uuid.UUID `gorm:"size:45;not null;column:principal_type_id" json:"principal_type_id"`
	Name            string    `gorm:"size:45;not null;column:name" json:"name"`
	RowStatus       int       `gorm:"default:1;column:row_status" json:"row_status"`
	CreatedBy       uuid.UUID `gorm:"size:45;column:created_by" json:"created_by"`
	UpdatedBy       uuid.UUID `gorm:"size:45;column:updated_by" json:"updated_by"`
	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at;autoCreateTime" json:"updated_at"`
}

func (t *MstPrincipalTypes) TableName() string {
	return "mst_principal_types"
}
