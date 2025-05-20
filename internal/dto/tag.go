package dto

type Tags struct {
	Key   string `json:"key" validate:"required,min=3,max=50"`
	Value string `json:"value" validate:"required,min=3,max=50"`
}
