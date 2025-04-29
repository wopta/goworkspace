package accounting

type Invoice interface {
	create(isPay bool, isProforma bool) (string, error)
	save(url string, path string) (string, error)
}

func DoInvoicePaid(inv Invoice, path string) (string, error) {
	url, err := inv.create(true, false)
	if err != nil {
		return "", err
	}
	return inv.save(url, path)
}
