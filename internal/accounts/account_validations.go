package accounts

import (
	"fmt"
	"iam_services_main_v1/gql/models"
	"iam_services_main_v1/internal/dto"
	"iam_services_main_v1/internal/validations"
	"iam_services_main_v1/pkg/logger"
	"log"
)

// ValidateCreateAccountInput validates the input for creating an account
func ValidateCreateAccountInput(input models.CreateAccountInput) (dto.CreateAccountInput, error) {
	// Map input to DTO
	dtoInput := MapToAccount(input)

	// Validate all components
	if err := validateComponents(dtoInput); err != nil {
		return dtoInput, err
	}

	return dtoInput, nil
}

// MapToAccount maps a TenantResource to an Account model
func MapToAccount(input models.CreateAccountInput) dto.CreateAccountInput {
	// Map tags from models.TagInput to dto.Tags
	tags := make([]dto.Tags, 0, len(input.Tags))
	for _, tag := range input.Tags {
		if tag != nil {
			tags = append(tags, dto.Tags{
				Key:   tag.Key,
				Value: tag.Value,
			})
		}
	}

	// Create the account DTO
	result := dto.CreateAccountInput{
		ID:             input.ID,
		Name:           input.Name,
		Description:    input.Description,
		TenantID:       input.TenantID,
		ParentID:       input.ParentID,
		RelationType:   string(input.RelationType),
		AccountOwnerID: input.AccountOwnerID,
		Tags:           tags,
	}

	// Only map billing info if it exists
	if input.BillingInfo != nil {
		result.BillingInfo = dto.BillingInfo{
			CreditCardNumber: input.BillingInfo.CreditCardNumber,
			CreditCardType:   input.BillingInfo.CreditCardType,
			ExpirationDate:   input.BillingInfo.ExpirationDate,
			Cvv:              input.BillingInfo.Cvv,
		}

		// Only map billing address if it exists
		if input.BillingInfo.BillingAddress != nil {
			result.BillingInfo.BillingAddress = dto.BillingAddress{
				Street:  input.BillingInfo.BillingAddress.Street,
				City:    input.BillingInfo.BillingAddress.City,
				State:   input.BillingInfo.BillingAddress.State,
				Country: input.BillingInfo.BillingAddress.Country,
				Zipcode: input.BillingInfo.BillingAddress.Zipcode,
			}
		}
	}

	return result
}

// validateComponents validates the components of the input
func validateComponents(input dto.CreateAccountInput) error {
	log.Println(input)
	// Only validate billing address if billing info exists
	if input.BillingInfo.CreditCardNumber != "" {
		log.Println("Validating billing info")
		// Validate billing info
		if err := validations.ValidateStruct(input.BillingInfo); err != nil {
			logger.LogError("billing info validation failed", "error", err)
			return fmt.Errorf("invalid billing info: %w", err)
		}

		// Validate billing address if it exists
		if input.BillingInfo.BillingAddress.Street != "" {
			log.Println("Validating billing info")
			if err := validations.ValidateStruct(input.BillingInfo.BillingAddress); err != nil {
				logger.LogError("billing address validation failed", "error", err)
				return fmt.Errorf("invalid billing address: %w", err)
			}
		}
	}

	// Validate entire input
	if err := validations.ValidateStruct(input); err != nil {
		logger.LogError("full input validation failed", "error", err)
		return fmt.Errorf("invalid input: %w", err)
	}

	return nil
}
