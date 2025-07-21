package dto

import "gitlab.dev.wopta.it/goworkspace/models"

type CatnatDTO struct {
}

func NewCatnatDto() CatnatDTO {
	return CatnatDTO{}
}
func (dto *CatnatDTO) FromPolicy(policy *models.Policy, node *models.NetworkNode) {
}
