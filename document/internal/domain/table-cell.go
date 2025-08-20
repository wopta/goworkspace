package domain

type TableCell struct {
	Text      string
	Height    float64
	Width     float64
	FontSize  FontSize
	FontStyle FontStyle
	FontColor Color
	Fill      bool
	FillColor Color
	Align     TextAlign
	Link      string
	Border    string
}
