package constants

import "gitlab.dev.wopta.it/goworkspace/document/internal/domain"

var (
	BlackColor = domain.Color{
		R: 0,
		G: 0,
		B: 0,
	}
	PinkColor = domain.Color{
		R: 229,
		G: 9,
		B: 117,
	}
	WhiteColor = domain.Color{
		R: 255,
		G: 255,
		B: 255,
	}
	GreyColor = domain.Color{
		R: 217,
		G: 217,
		B: 217,
	}
	LightGreyColor = domain.Color{
		R: 242,
		G: 242,
		B: 242,
	}
	WatermarkColor = domain.Color{
		R: 206,
		G: 216,
		B: 232,
	}
	YellowColor = domain.Color{
		R: 255,
		G: 255,
		B: 0,
	}
	NoColor = domain.Color{}
)
