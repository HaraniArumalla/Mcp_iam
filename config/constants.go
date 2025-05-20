package config

// Database configuration variables
type contextKey string

const GinContextKey contextKey = "GinContextKey"

const (
	DBLoc               = "Local"
	DBParseTime         = "True"
	DBCharset           = "utf8mb4"
	AccountsContextKey  = "accounts"
	GenericErrorMessage = "An error has occurred. Please try again later or contact support."
)

// Permit configuration variables

const (
	ClientOrgUnitResourceTypeID = "ed113dd2-bbda-11ef-87ea-c03c5946f955"
	AccountResourceTypeID       = "ed113f30-bbda-11ef-87ea-c03c5946f955"
	TenantResourceTypeID        = "ed113bda-bbda-11ef-87ea-c03c5946f955"
	RoleResourceTypeID          = "464b359e-3d43-4461-bb92-d36ebaf29082"
)

// Constant configuration variables
const (
	Tenant                 = "Tenant"
	ClientOrganizationUnit = "ClientOrganizationUnit"
	Account                = "Account"
	Role                   = "Role"
)
