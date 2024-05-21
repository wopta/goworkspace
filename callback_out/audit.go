package callback_out

import (
	"io"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/civil"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

const CallbackOutTableId string = "callback-out"

func saveAudit(node *models.NetworkNode, action CallbackoutAction, req *http.Request, res *http.Response) {
	var audit struct {
		creationDate civil.DateTime
		client       string
		nodeUid      string
		action       string
		reqMethod    string
		reqPath      string
		reqBody      string
		resStatus    string
		resBody      string
	}

	reqBody, _ := io.ReadAll(req.Body)
	resBody, _ := io.ReadAll(res.Body)
	defer func() {
		req.Body.Close()
		res.Body.Close()
	}()

	audit.creationDate = civil.DateTimeOf(time.Now().UTC())
	audit.client = node.CallbackConfig.Name
	audit.nodeUid = node.Uid
	audit.action = action
	audit.reqMethod = req.Method
	audit.reqPath = req.URL.RequestURI()
	audit.reqBody = string(reqBody)
	audit.resStatus = res.Status
	audit.resBody = string(resBody)

	if err := lib.InsertRowsBigQuery(lib.WoptaDataset, CallbackOutTableId, audit); err != nil {
		log.Printf("error saving audit: %s", err)
	}
}
