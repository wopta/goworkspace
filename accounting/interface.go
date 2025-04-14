package accounting

type Invoice interface {
	Create(isPay bool, isProforma bool) string

	Save(url string, path string) error
}

func DoInvicePaid(inv Invoice, path string) {

	url := inv.Create(true, false)
	inv.Save(url, path)

}
