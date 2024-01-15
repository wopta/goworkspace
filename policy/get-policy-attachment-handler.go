package policy

import (
	"encoding/base64"
	"encoding/json"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type GetAttachmentReq struct {
	PolicyUid string `json:"policyUid"`
	Filename  string `json:"filename"`
}

type GetAttachmentResp struct {
	RawDoc string `json:"rawDoc"`
}

func GetAttachmentFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err    error
		policy models.Policy
		req    GetAttachmentReq
		resp   GetAttachmentResp
	)

	log.SetPrefix("[GetAttachmentFx]")

	log.Println("Handler start -----------------------------------------------")

	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()
	log.Printf("request body: %s", string(body))
	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Printf("error unmarshaling request: %s", err.Error())
		return "", nil, err
	}

	log.Printf("retrieving policy %s from Firestore...", req.PolicyUid)

	policy, err = GetPolicy(req.PolicyUid, "")
	if err != nil {
		log.Printf("error retrieving policy %s from Firestore: %s", req.PolicyUid, err.Error())
		return "", nil, err
	}

	log.Println("checking if requested attachment is present...")

	if policy.Attachments != nil {
		for _, att := range *policy.Attachments {
			if strings.EqualFold(att.FileName, req.Filename) {
				rawDoc, err := downloadAttachment(att.Link)
				if err != nil {
					log.Printf("error downloading attacchment %s from Google Bucket: %s", att.FileName, err.Error())
					return "", nil, err
				}
				resp = GetAttachmentResp{
					RawDoc: rawDoc,
				}
				rawResp, err := json.Marshal(resp)
				return string(rawResp), rawResp, err
			}
		}
	}
	return "", nil, err
}

func downloadAttachment(gsLink string) (string, error) {
	if !strings.Contains(gsLink, "gs://") {
		gsLink = "gs://" + os.Getenv("GOOGLE_STORAGE_BUCKET") + "/" + gsLink
	}
	rawDoc, err := lib.ReadFileFromGoogleStorage(gsLink)
	if err != nil {
		log.Printf("[GetPolicyFx] error reading document from Google Storage: %s", err.Error())
		return "", err
	}
	return base64.StdEncoding.EncodeToString(rawDoc), nil
}
