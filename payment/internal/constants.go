package internal

import "errors"

var (
	ErrTransactionPaid         = errors.New("transaction already paid")
	ErrTransactionOutOfOrder   = errors.New("transaction is not next on schedule. Payment must be ordered")
	ErrPolicyNotSigned         = errors.New("policy not signed")
	ErrPaymentMethodNotAllowed = errors.New("payment method not allowed")
	ErrPaymentFailed           = errors.New("failed to update DB")
	ErrTransactionDeleted      = errors.New("transaction already deleted")
	ErrInvalidTransactions     = errors.New("invalid number of transactions")
)
