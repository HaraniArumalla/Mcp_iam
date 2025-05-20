package dto

import (
	"github.com/google/uuid"
)

type BillingAddress struct {
	Street  string `json:"street" validate:"required,min=5,max=100"`
	City    string `json:"city" validate:"required,min=2,max=50"`
	State   string `json:"state" validate:"required,min=2,max=50"`
	Country string `json:"country" validate:"required,iso3166_1_alpha2"`
	Zipcode string `json:"zipcode" validate:"required,numeric,min=5,max=10"`
}

type BillingInfo struct {
	CreditCardNumber string         `json:"creditCardNumber" validate:"required,numeric,len=16"`
	CreditCardType   string         `json:"creditCardType" validate:"required,oneof=visa mastercard amex discover"`
	ExpirationDate   string         `json:"expirationDate" validate:"required,datetime=01/06"`
	Cvv              string         `json:"cvv" validate:"required,numeric,len=3"`
	BillingAddress   BillingAddress `json:"billingAddress" validate:"omitempty"`
}

type CreateAccountInput struct {
	ID             uuid.UUID   `json:"id" validate:"required"`
	Name           string      `json:"name" validate:"required,min=3,max=100"`
	Description    *string     `json:"description" validate:"omitempty,max=500"`
	TenantID       uuid.UUID   `json:"tenantId" validate:"required"`
	ParentID       uuid.UUID   `json:"parentId" validate:"required"`
	RelationType   string      `json:"relationType" validate:"required,oneof=CHILD PARENT SELF"`
	AccountOwnerID uuid.UUID   `json:"accountOwnerId" validate:"required"`
	BillingInfo    BillingInfo `json:"billingInfo" validate:"omitempty"`
	Tags           []Tags      `json:"tags" validate:"omitempty"`
}
