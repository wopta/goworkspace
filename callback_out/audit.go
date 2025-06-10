package callback_out

import (
	"time"

	"cloud.google.com/go/bigquery"
	"gitlab.dev.wopta.it/goworkspace/callback_out/base"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models"
)

const CallbackOutTableId string = "callback-out"

type auditSchema struct {
	CreationDate  bigquery.NullDateTime `bigquery:"creationDate"`
	Client        string                `bigquery:"client"`
	NodeUid       string                `bigquery:"nodeUid"`
	Action        string                `bigquery:"action"`
	ReqMethod     string                `bigquery:"reqMethod"`
	ReqPath       string                `bigquery:"reqPath"`
	ReqBody       string                `bigquery:"reqBody"`
	ResStatusCode int                   `bigquery:"resStatusCode"`
	ResBody       string                `bigquery:"resBody"`
	Error         string                `bigquery:"error"`
}

func saveAudit(node *models.NetworkNode, action base.CallbackoutAction, res base.CallbackInfo) {
	var (
		audit auditSchema
	)

	audit.CreationDate = lib.GetBigQueryNullDateTime(time.Now().UTC())
	audit.Client = node.CallbackConfig.Name
	audit.NodeUid = node.Uid
	audit.Action = string(action)

	audit.ReqBody = string(res.ReqBody)
	audit.ReqMethod = res.ReqMethod
	audit.ReqPath = res.ReqPath
	//audit.ReqPath = res.Request.Host + res.Request.URL.RequestURI()

	audit.ResStatusCode = res.ResStatusCode
	audit.ResBody = string(res.ResBody)

	if res.Error != nil {
		audit.Error = res.Error.Error()
	}

	if err := lib.InsertRowsBigQuery(lib.WoptaDataset, CallbackOutTableId, audit); err != nil {
		log.ErrorF("error saving audit: %s", err)
	}
}
