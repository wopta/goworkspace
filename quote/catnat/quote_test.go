package catnat

import (
	"errors"
	"testing"

	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/sellable"
)

func mock_sellableWithError(_ *models.Policy, _ *models.Product, _ bool) (*sellable.SellableOutput, error) {
	return &sellable.SellableOutput{}, errors.New("sellable went bad")
}

func mock_sellable(_ *models.Policy, _ *models.Product, _ bool) (*sellable.SellableOutput, error) {
	return &sellable.SellableOutput{
		Product: &models.Product{},
	}, nil
}

type mock_clientCatnat struct {
	withError bool
}

func (c *mock_clientCatnat) Quote(dto RequestDTO) (response ResponseDTO, err error) {
	if c.withError {
		return response, errors.New("quote error")
	}
	return response, nil
}
func (c *mock_clientCatnat) Emit(dto RequestDTO) (response any, err error) {
	return response, nil
}

func TestQuoteCatnat(t *testing.T) {
	_, err := CatnatQuote(new(models.Policy), new(models.Product), mock_sellable, new(mock_clientCatnat))
	if err != nil {
		t.Fatal(err)
	}
}
