package policy

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type UploadPolicyMediaReq struct {
	PolicyUid string `json:"policyUid"`
	Filename  string `json:"filename"`
	Base64    string `json:"base64"`
	MimeType  string `json:"mimeType"`
	Name      string `json:"name"`
	Section   string `json:"section"`
	IsPrivate bool   `json:"isPrivate"`
	Note      string `json:"note"`
}

func UploadPolicyMediaFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err    error
		policy *models.Policy
		req    UploadPolicyMediaReq
	)

	log.SetPrefix("[UploadPolicyMediaFx]")

	log.Println("Handler start -----------------------------------------------")

	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()
	log.Printf("request body: %s", string(body))
	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Printf("error unmarshaling request: %s", err.Error())
		return "", nil, err
	}

	now := time.Now().UTC()
	req.Filename = fmt.Sprintf("%s_%s", req.Filename, now)

	log.Printf("getting policy %s from Firestore...")

	docSnap, err := lib.GetFirestoreErr(models.PolicyCollection, req.PolicyUid)
	if err != nil {
		log.Printf("error getting policy %s from Firestore: %s", req.PolicyUid, err.Error())
		return "", nil, err
	}
	err = docSnap.DataTo(&policy)
	if err != nil {
		log.Printf("error converting docsnap to policy: %s", err.Error())
		return "", nil, err
	}

	err = putAttachment(policy, req)

	log.Println("Handler end -------------------------------------------------")
	log.SetPrefix("")

	return "", nil, err
}

func putAttachment(policy *models.Policy, req UploadPolicyMediaReq) error {
	log.Printf("converting base64 to []byte")
	rawDoc, err := base64.StdEncoding.DecodeString(req.Base64)

	log.Printf("uploading %s to asset/users/%s in Google Bucket", req.Filename, policy.Contractor.Uid)

	gsLink, err := lib.PutToGoogleStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), fmt.Sprintf("assets/users/%s/%s",
		policy.Contractor.Uid, req.Filename), rawDoc)
	if err != nil {
		log.Printf("error uploading %s to Google Bucket: %s", req.Filename, err.Error())
		return err
	}

	att := models.Attachment{
		Name:      strings.Split(req.Name, ".")[0],
		Link:      gsLink,
		FileName:  req.Filename,
		MimeType:  req.MimeType,
		IsPrivate: req.IsPrivate,
		Section:   req.Section,
		Note:      req.Note,
	}

	if policy.Attachments == nil {
		policy.Attachments = new([]models.Attachment)
	}
	*policy.Attachments = append(*policy.Attachments, att)

	log.Printf("saving policy %s to Firestore...", policy.Uid)

	err = lib.SetFirestoreErr(models.PolicyCollection, policy.Uid, policy)
	if err != nil {
		log.Printf("error saving policy %s to Firestore: %s", policy.Uid, err.Error())
		return err
	}
	log.Printf("policy %s saved into Firestore", policy.Uid)

	log.Printf("saving policy %s to BigQuery...", policy.Uid)

	policy.BigquerySave("")

	return err
}
