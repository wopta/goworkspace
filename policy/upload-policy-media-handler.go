package policy

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

type UploadPolicyMediaReq struct {
	PolicyUid string
	Filename  string
	Bytes     []byte
	MimeType  string
	Name      string
	Section   string
	IsPrivate bool
	Note      string
}

func UploadPolicyMediaFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err    error
		policy *models.Policy
		req    UploadPolicyMediaReq
	)

	log.SetPrefix("[UploadPolicyMediaFx]")
	defer log.SetPrefix("")

	log.Println("Handler start -----------------------------------------------")

	err = r.ParseMultipartForm(32 << 20)
	if err != nil {
		log.Printf("error parsing multipart form: %s", err.Error())
		return "", nil, err
	}
	req.PolicyUid = r.PostFormValue("policyUid")
	req.Filename = r.PostFormValue("filename")
	req.MimeType = r.PostFormValue("mimeType")
	req.Name = r.PostFormValue("name")
	req.Section = r.PostFormValue("section")
	req.IsPrivate = r.PostFormValue("isPrivate") == "true"
	req.Note = r.PostFormValue("note")
	file, _, err := r.FormFile("bytes")

	if err != nil {
		log.Printf("error getting file from request: %s", err.Error())
		return "", nil, err
	}

	log.Printf("policyUid: %s", req.PolicyUid)
	log.Printf("filename: %s", req.Filename)
	log.Printf("mimeType: %s", req.MimeType)
	log.Printf("name: %s", req.Name)
	log.Printf("section: %s", req.Section)
	log.Printf("isPrivate: %t", req.IsPrivate)
	log.Printf("note: %s", req.Note)

	defer file.Close()
	req.Bytes, err = io.ReadAll(file)
	if err != nil {
		log.Printf("error reading file from request: %s", err.Error())
		return "", nil, err
	}

	log.Printf("getting policy %s from Firestore...", req.PolicyUid)

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

	return "", nil, err
}

func putAttachment(policy *models.Policy, req UploadPolicyMediaReq) error {
	var filename string

	splittedFilename := strings.Split(req.Filename, ".")
	if len(splittedFilename) > 2 {
		filename = strings.Join(splittedFilename[:len(splittedFilename)-1], ".")
	} else {
		filename = splittedFilename[0]
	}
	filename += fmt.Sprintf("_%d.%s", time.Now().UTC().Unix(), splittedFilename[len(splittedFilename)-1])

	log.Printf("uploading %s to asset/users/%s in Google Bucket", filename, policy.Contractor.Uid)

	gsLink, err := lib.PutToGoogleStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), fmt.Sprintf("assets/users/%s/%s",
		policy.Contractor.Uid, filename), req.Bytes)
	if err != nil {
		log.Printf("error uploading %s to Google Bucket: %s", filename, err.Error())
		return err
	}

	att := models.Attachment{
		Name:      req.Name,
		Link:      gsLink,
		FileName:  filename,
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
