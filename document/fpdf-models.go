package document

type rgbColor struct {
	r int
	g int
	b int
}

type tableCell struct {
	text      string
	height    float64
	width     float64
	textBold  bool
	fill      bool
	fillColor rgbColor
	align     string
	border    string
}
