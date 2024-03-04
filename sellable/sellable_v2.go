package sellable

import (
	"errors"
	"log"

	"github.com/wopta/goworkspace/models"
)

func SellableV2(inputPolicy models.Policy) (product models.Product) {
	// validate input
	// // life -> contractor age
	// // persona -> ?
	// // pmi -> ?
	// // gap -> policyDuration X vehicle.Registration

	// extract desired fields from input
	// // life -> contractor age
	// // persona -> ?
	// // pmi -> ?
	// // gap -> policyDuration X vehicle.Registration
	getInputByProduct(inputPolicy)

	// load base product
	// // get from bucket
	// // get from fs

	// compute product
	// // by rules engine
	// // by code

	// return computed
	return product
}

type SellableInput interface {
	GetInput(policy models.Policy) interface{}
	GetReturnType() string
}

type LifeSellableInput struct {
	returnType string
}

func (lsi LifeSellableInput) GetInput(policy models.Policy) interface{} {
	return map[string]int{
		"age": 18,
	}
}

func (lsi LifeSellableInput) GetReturnType() string {
	return lsi.returnType
}

func getInputByProduct(policy models.Policy) (SellableInput, error) {
	switch policy.Name {
	case models.LifeProduct:
		return LifeSellableInput{returnType: "map[string]int"}, nil
	case models.PersonaProduct:
		return nil, errors.New("not implemented")
	case models.PmiProduct:
		return nil, errors.New("not implemented")
	case models.GapProduct:
		return nil, errors.New("not implemented")
	default:
		log.Printf("product %s not implemented", policy.Name)
		return nil, errors.New("not implemented")
	}
}
