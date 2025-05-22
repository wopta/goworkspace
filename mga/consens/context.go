package consens

import (
	"context"
	"time"

	"gitlab.dev.wopta.it/goworkspace/lib/log"
)

type key string

const timestamp = key("timestamp")

func getTimestamp(ctx context.Context) time.Time {
	if rawTime := ctx.Value(timestamp); rawTime != nil {
		return (rawTime).(time.Time)
	}
	log.Println("timestamp not set - defaulting to time.Now().UTC()")
	return time.Now().UTC()
}
