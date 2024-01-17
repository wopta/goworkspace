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

type GetPolicyMediaReq struct {
	PolicyUid string `json:"policyUid"`
	Filename  string `json:"filename"`
}

type GetPolicyMediaResp struct {
	RawDoc string `json:"rawDoc"`
}

func GetPolicyMediaFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err    error
		policy models.Policy
		req    GetPolicyMediaReq
	)

	log.SetPrefix("[GetPolicyMediaFx]")

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
				rawResp, resp, err := downloadAttachment(att.Link)
				if err != nil {
					log.Printf("error downloading media %s from Google Bucket: %s", req.Filename, err.Error())
					return "", nil, err
				}
				return rawResp, resp, err
			}
		}
	}

	if policy.ReservedInfo != nil {
		for _, doc := range policy.ReservedInfo.Documents {
			if strings.EqualFold(doc.FileName, req.Filename) {
				rawResp, resp, err := downloadAttachment(doc.Link)
				if err != nil {
					log.Printf("error downloading media %s from Google Bucket: %s", req.Filename, err.Error())
					return "", nil, err
				}
				return rawResp, resp, err
			}
		}
	}

	for _, doc := range policy.Contractor.IdentityDocuments {
		if strings.EqualFold(doc.FrontMedia.FileName, req.Filename) {
			rawResp, resp, err := downloadAttachment(doc.FrontMedia.Link)
			if err != nil {
				log.Printf("error downloading media %s from Google Bucket: %s", req.Filename, err.Error())
				return "", nil, err
			}
			return rawResp, resp, err
		}

		if doc.BackMedia != nil && strings.EqualFold(doc.BackMedia.FileName, req.Filename) {
			rawResp, resp, err := downloadAttachment(doc.BackMedia.Link)
			if err != nil {
				log.Printf("error downloading media %s from Google Bucket: %s", req.Filename, err.Error())
				return "", nil, err
			}

			return rawResp, resp, err
		}
	}

	log.Println("Handler end -------------------------------------------------")

	return "", nil, err
}

func downloadAttachment(gsLink string) (string, GetPolicyMediaResp, error) {
	if !strings.Contains(gsLink, "gs://") {
		gsLink = "gs://" + os.Getenv("GOOGLE_STORAGE_BUCKET") + "/" + gsLink
	}
	rawDoc, err := lib.ReadFileFromGoogleStorage(gsLink)
	if err != nil {
		log.Printf("error reading document from Google Storage: %s", err.Error())
		return "", GetPolicyMediaResp{}, err
	}
	log.Printf("document found")

	resp := GetPolicyMediaResp{
		RawDoc: base64.StdEncoding.EncodeToString(rawDoc),
	}

	rawResp, err := json.Marshal(resp)
	return string(rawResp), resp, err
}
