package dto

import (
	"github.com/google/uuid"
)

type Address struct {
	Street  string `json:"street" validate:"required,min=5,max=100"`
	City    string `json:"city" validate:"required,min=2,max=50"`
	State   string `json:"state" validate:"required,min=2,max=50"`
	Country string `json:"country" validate:"required"`
	Zipcode string `json:"zipcode" validate:"required,numeric,min=5,max=10"`
}

type ContactInfo struct {
	Email       string  `json:"email" validate:"required,email"`
	PhoneNumber string  `json:"phoneNumber" validate:"required,numeric,min=10,max=15"`
	Address     Address `json:"Address" validate:"required"`
}

type CreateTenantInput struct {
	ID             uuid.UUID   `json:"id" validate:"required"`
	Name           string      `json:"name" validate:"required,min=3,max=100"`
	Description    *string     `json:"description" validate:"omitempty,max=500"`
	ContactInfo    ContactInfo `json:"contactInfo" validate:"required"`
	Tags           []Tags      `json:"tags" validate:"omitempty"`
	AccountOwnerId uuid.UUID   `json:"accountOwnerId" validate:"required"`
}
