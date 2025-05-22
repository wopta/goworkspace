package dto

import (
	"fmt"

	"gitlab.dev.wopta.it/goworkspace/document/internal/constants"
	"gitlab.dev.wopta.it/goworkspace/lib"
)

type claimDTO struct {
	Description string
	Quantity    numeric
	Value       numeric
}

func newClaimDTO() *claimDTO {
	return &claimDTO{
		Description: constants.EmptyField,
		Quantity: numeric{
			ValueFloat: 0,
			ValueInt:   0,
			Text:       fmt.Sprintf("%d", 0),
		},
		Value: numeric{
			ValueFloat: 0,
			ValueInt:   0,
			Text:       lib.HumanaizePriceEuro(0),
		},
	}
}
