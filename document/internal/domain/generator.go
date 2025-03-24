package domain

type QuoteGenerator interface {
	Exec() ([]byte, error)
}
