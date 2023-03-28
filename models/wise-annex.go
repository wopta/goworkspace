package models

import (
	"encoding/json"
	"fmt"
	wiseProxy "github.com/wopta/goworkspace/wiseproxy"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type WiseAnnex struct {
	Id         string    `json:"txRifIdAllegato,omitempty"`
	Name       string    `json:"txNomeAllegato,omitempty"`
	InsertDate time.Time `json:"dtInserimento,omitempty"`
}

func (annex WiseAnnex) ToDomain(wiseToken *string) (Attachment, *string) {
	var (
		responseReader io.ReadCloser
		wiseAnnex      WiseBase64Annex
	)

	request := []byte(fmt.Sprintf(`{"txRifAllegato": "%s", "cdLingua": "it"}`, annex.Id))
	responseReader, wiseToken = wiseProxy.WiseBatch("WebApiProduct/Api/recuperaAllegato", request, http.MethodPost, wiseToken)

	jsonData, _ := ioutil.ReadAll(responseReader)

	_ = json.Unmarshal(jsonData, &wiseAnnex)

	var attachment Attachment

	attachment.Name = annex.Name
	attachment.Byte = wiseAnnex.Bytes

	return attachment, wiseToken
}
