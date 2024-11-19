package manual

import "errors"

var (
	errTransactionPaid         = errors.New("transaction already paid")
	errTransactionOutOfOrder   = errors.New("transaction is not next on schedule. Payment must be ordered")
	errPolicyNotSigned         = errors.New("policy not signed")
	errPaymentMethodNotAllowed = errors.New("payment method not allowed")
	errPaymentFailed           = errors.New("failed to update DB")
	errTransactionDeleted      = errors.New("transaction already deleted")
)
