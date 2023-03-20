package lib

import (
	"time"
)

func Dateformat(t time.Time) string {
	layout := "02/01/2006"
	return t.Format(layout)
}
