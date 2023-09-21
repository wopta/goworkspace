package product

import (
	"github.com/wopta/goworkspace/models"
)

var (
	ageMap = map[string]map[string]map[string]int{
		models.UserRoleAdmin: {
			models.LifeProduct: {
				minAge:         75,
				minReservedAge: 55,
			},
			models.PersonaProduct: {
				minAge:         75,
				minReservedAge: 75,
			},
		},
		models.UserRoleManager: {
			models.LifeProduct: {
				minAge:         75,
				minReservedAge: 55,
			},
			models.PersonaProduct: {
				minAge:         75,
				minReservedAge: 75,
			},
		},
		models.UserRoleAgent: {
			models.LifeProduct: {
				minAge:         70,
				minReservedAge: 55,
			},
			models.PersonaProduct: {
				minAge:         75,
				minReservedAge: 75,
			},
		},
		models.UserRoleAgency: {
			models.LifeProduct: {
				minAge:         71,
				minReservedAge: 55,
			},
			models.PersonaProduct: {
				minAge:         75,
				minReservedAge: 75,
			},
		},
		models.UserRoleAll: {
			models.LifeProduct: {
				minAge:         55,
				minReservedAge: 55,
			},
			models.PersonaProduct: {
				minAge:         75,
				minReservedAge: 75,
			},
		},
		models.UserRoleCustomer: {
			models.LifeProduct: {
				minAge:         55,
				minReservedAge: 55,
			},
			models.PersonaProduct: {
				minAge:         75,
				minReservedAge: 75,
			},
		},
	}
)

func GetReservedAge(product, channel string) int {
	ret := 0
	if ageMap[channel] != nil && ageMap[channel][product] != nil {
		ret = ageMap[channel][product][minReservedAge]
	}
	return ret
}
