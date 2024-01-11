package payment

const (
	errTransactionPaid         = "MPTR001: Transaction already paid"
	errTransactionOutOfOrder   = "MPTR002: Transaction is not next on schedule. Payment must be ordered"
	errPolicyNotSigned         = "MPTR003: Policy not signed"
	errPaymentMethodNotAllowed = "MPTR004: Payment method not allowed"
	errPaymentFailed           = "MPTR005: Failed to update DB"
	errTransactionDeleted      = "MPTR006: Transaction already deleted"
)
