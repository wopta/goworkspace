package callbackout

import (
	"errors"

	"github.com/wopta/goworkspace/models"
)

var handlerMap map[string]func(models.Policy) error = map[string]func(models.Policy) error{
	"winCallbackHandler": winCallbackHandler,
}
var ErrCallbackNotSet = errors.New("callback not set")

func Handler(fxName string, policy models.Policy) error {
	var (
		callback func(models.Policy) error
		ok       bool
	)

	if callback, ok = handlerMap[fxName]; !ok {
		return ErrCallbackNotSet
	}

	return callback(policy)
}
