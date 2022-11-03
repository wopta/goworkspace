package document

import "log"

func (s Skin) lenToHeight(w string) float64 {

	if len(w) > s.DynamicHeightMin {
		log.Println((float64(len(w)) / 30.0))
		return (float64(len(w)) / s.DynamicHeightDiv)
	} else {
		return s.rowHeight
	}

}
