package tenants

import (
	"fmt"
	"iam_services_main_v1/gql/models"
	"iam_services_main_v1/internal/dto"
	"iam_services_main_v1/internal/validations"
	"iam_services_main_v1/pkg/logger"
)

// ValidateCreateTenantInput validates the input for creating an tenant
func ValidateCreateTenantInput(input models.CreateTenantInput) (dto.CreateTenantInput, error) {
	// Early validation of required fields
	if err := validateRequiredFields(input); err != nil {
		return dto.CreateTenantInput{}, err
	}

	// Map input to DTO
	dtoInput := MapToTenant(input)

	// Validate all components
	if err := validateComponents(dtoInput); err != nil {
		return dtoInput, err
	}

	return dtoInput, nil
}

// validateRequiredFields validates the required fields of the input
func validateRequiredFields(input models.CreateTenantInput) error {
	if input.ContactInfo == nil {
		return fmt.Errorf("contact info is required")
	}
	if input.ContactInfo.Address == nil {
		return fmt.Errorf("address is required")
	}
	return nil
}

// MapToTenant maps a TenantResource to an Tenant model
func MapToTenant(input models.CreateTenantInput) dto.CreateTenantInput {
	dtoInput := dto.CreateTenantInput{
		ID:             input.ID,
		Name:           input.Name,
		Description:    input.Description,
		AccountOwnerId: input.AccountOwnerID,
		ContactInfo: dto.ContactInfo{
			Email:       *input.ContactInfo.Email,
			PhoneNumber: *input.ContactInfo.PhoneNumber,
			Address: dto.Address{
				Street:  *input.ContactInfo.Address.Street,
				City:    *input.ContactInfo.Address.City,
				State:   *input.ContactInfo.Address.State,
				Country: *input.ContactInfo.Address.Country,
				Zipcode: *input.ContactInfo.Address.Zipcode,
			},
		},
	}

	// Add Tags mapping if input has tags
	if len(input.Tags) > 0 {
		dtoInput.Tags = make([]dto.Tags, len(input.Tags))
		for i, tag := range input.Tags {
			dtoInput.Tags[i] = dto.Tags{
				Key:   tag.Key,
				Value: tag.Value,
			}
		}
	}
	logger.LogInfo("Mapped input to tenant DTO", "input", input, "dtoInput", dtoInput)
	return dtoInput
}

// function validateComponents validates the components of the input
func validateComponents(input dto.CreateTenantInput) error {
	// Validate address
	if err := validations.ValidateStruct(input.ContactInfo.Address); err != nil {
		logger.LogError("billing address validation failed", "error", err)
		return fmt.Errorf("invalid address: %w", err)
	}

	// Validate contact info
	if err := validations.ValidateStruct(input.ContactInfo); err != nil {
		logger.LogError("contact info validation failed", "error", err)
		return fmt.Errorf("invalid contact info: %w", err)
	}

	// Validate entire input
	if err := validations.ValidateStruct(input); err != nil {
		logger.LogError("full input validation failed", "error", err)
		return fmt.Errorf("invalid input: %w", err)
	}

	return nil
}
