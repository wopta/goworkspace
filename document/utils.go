package document

import "log"

func (s Skin) lenToHeight(w string) float64 {

	if len(w) > 90 {
		log.Println((float64(len(w)) / 30.0))
		return (float64(len(w)) / 32.0)
	} else {
		return s.rowHeight
	}

}
