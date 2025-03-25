package accounting

type Invoice interface {
	Create(isPay bool, isProforma bool)
	Sand()
}
