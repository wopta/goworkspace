package callback_out

import (
	"io"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

const CallbackOutTableId string = "callback-out"

type auditSchema struct {
	CreationDate bigquery.NullDateTime `bigquery:"creationDate"`
	Client       string                `bigquery:"client"`
	NodeUid      string                `bigquery:"nodeUid"`
	Action       string                `bigquery:"action"`
	ReqMethod    string                `bigquery:"reqMethod"`
	ReqPath      string                `bigquery:"reqPath"`
	ReqBody      string                `bigquery:"reqBody"`
	ResStatus    string                `bigquery:"resStatus"`
	ResBody      string                `bigquery:"resBody"`
}

func saveAudit(node *models.NetworkNode, action CallbackoutAction, req *http.Request, res *http.Response) {
	var audit auditSchema

	reqBody, _ := io.ReadAll(req.Body)
	resBody, _ := io.ReadAll(res.Body)
	defer func() {
		req.Body.Close()
		res.Body.Close()
	}()

	audit.CreationDate = lib.GetBigQueryNullDateTime(time.Now().UTC())
	audit.Client = node.CallbackConfig.Name
	audit.NodeUid = node.Uid
	audit.Action = action
	audit.ReqMethod = req.Method
	audit.ReqPath = req.URL.RequestURI()
	audit.ReqBody = string(reqBody)
	audit.ResStatus = res.Status
	audit.ResBody = string(resBody)

	if err := lib.InsertRowsBigQuery(lib.WoptaDataset, CallbackOutTableId, audit); err != nil {
		log.Printf("error saving audit: %s", err)
	}
}
