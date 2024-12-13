package domain

type TableCell struct {
	Text      string
	Height    float64
	Width     float64
	TextBold  bool
	Fill      bool
	FillColor Color
	Align     TextAlign
	Border    string
}
