package common

import "time"

type GitTag struct {
	Name      string
	Module    string
	Version   string
	Env       string
	Hash      string
	CreatedAt time.Time
}
