package callback_out

import (
	"io"
	"log"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/wopta/goworkspace/callback_out/internal"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
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

func saveAudit(node *models.NetworkNode, action CallbackoutAction, res internal.CallbackInfo) {
	var audit auditSchema

	resBody, _ := io.ReadAll(res.Response.Body)
	defer res.Response.Body.Close()

	audit.CreationDate = lib.GetBigQueryNullDateTime(time.Now().UTC())
	audit.Client = node.CallbackConfig.Name
	audit.NodeUid = node.Uid
	audit.Action = action
	audit.ReqMethod = res.Request.Method
	audit.ReqPath = res.Request.Host + res.Request.URL.RequestURI()
	audit.ReqBody = string(res.RequestBody)
	audit.ResStatusCode = res.Response.StatusCode
	audit.ResBody = string(resBody)
	if res.Error != nil {
		audit.Error = res.Error.Error()
	}

	if err := lib.InsertRowsBigQuery(lib.WoptaDataset, CallbackOutTableId, audit); err != nil {
		log.Printf("error saving audit: %s", err)
	}
}
