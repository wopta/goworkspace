package payment

const (
	ERR_TR_PAID         = "MPTR001: Transaction already paid."
	ERR_TR_OUT_OF_ORDER = "MPTR002: Transaction is not next on schedule. Payment must be ordered."
	ERR_PL_NOT_SIGNED   = "MPTR003: Policy not signed"
)
