package domain

type TableCell struct {
	Text      string
	Height    float64
	Width     float64
	FontSize  FontSize
	FontStyle FontStyle
	Fill      bool
	FillColor Color
	Align     TextAlign
	Border    string
}
