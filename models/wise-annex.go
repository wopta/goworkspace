package models

import (
	"time"
)

type WiseAnnex struct {
	Id         string    `json:"txRifIdAllegato,omitempty"`
	Name       string    `json:"txNomeAllegato,omitempty"`
	InsertDate time.Time `json:"dtInserimento,omitempty"`
}

func (annex WiseAnnex) ToDomain(wiseToken *string) (Attachment, *string) {
	var attachment Attachment

	attachment.Name = annex.Name

	return attachment, wiseToken
}
