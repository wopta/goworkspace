package payment

const (
	errTransactionPaid         = "MPTR001: Transaction already paid"
	errTransactionOutOfOrder   = "MPTR002: Transaction is not next on schedule. Payment must be ordered"
	errPolicyNotSigned         = "MPTR003: Policy not signed"
	errPaymentMethodNotAllowed = "MPTR004: Payment method not allowed"
)

const (
	payMethodCard       = "creditcard"
	payMethodTransfer   = "transfer"
	payMethodSdd        = "sdd"
	PayMethodRemittance = "remittance"
)

func GetAllPaymentMethods() []string {
	return []string{payMethodCard, payMethodTransfer, payMethodSdd}
}
