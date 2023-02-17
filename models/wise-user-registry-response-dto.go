package models

type WiseUserRegistryResponseDto struct {
	UserRegistries *[]WiseUserAddressRegistryDto `json:"listAnagrafiche,omitempty"`
}
