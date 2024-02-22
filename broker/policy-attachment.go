package broker

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	plc "github.com/wopta/goworkspace/policy"
	"github.com/wopta/goworkspace/wiseproxy"
)

func GetPolicyAttachmentFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	var response GetPolicyAttachmentsResponse

	attachments, err := GetPolicyAttachments(r.Header.Get("policyUid"), r.Header.Get("origin"))
	if err != nil {
		log.Println("GetPolicyAttachments Error: " + err.Error())
		return "{}", nil, err
	}

	response.Attachments = attachments
	res, err := json.Marshal(response)
	if err != nil {
		log.Println("AttachmentsMarshal Error: " + err.Error())
	}

	return string(res), nil, err
}

func GetPolicyAttachments(policyUid string, origin string) ([]models.Attachment, error) {
	var (
		err         error
		wiseToken   *string = nil
		attachments []models.Attachment
	)
	if strings.HasPrefix(policyUid, "wise:") {
		// get attachment from wise
		var (
			wisePolicy         WiseCompletePolicyResponse
			attachmentResponse WiseAttachmentResponse
			wiseResponseData   []byte
		)
		wiseAttachmentId := strings.Split(policyUid, ":")[1]
		request := []byte(fmt.Sprintf(`{"idPolizza": "%s", "cdLingua": "it"}`, wiseAttachmentId))
		ioReader := wiseproxy.WiseProxyObj("WebApiProduct/Api/GetPolizzaCompleta", request, http.MethodPost)

		defer ioReader.Close()
		wiseResponseData, err = io.ReadAll(ioReader)

		if err != nil {
			return make([]models.Attachment, 0), err
		}
		err = json.Unmarshal(wiseResponseData, &wisePolicy)

		for _, wiseAttachment := range wisePolicy.Policy.Attachments {
			var attachment models.Attachment

			request := []byte(fmt.Sprintf(`{"txRifAllegato": "%s", "cdLingua": "it"}`, wiseAttachment.Id))
			ioReader, wiseToken = wiseproxy.WiseBatch("WebApiProduct/Api/recuperaAllegato", request, http.MethodPost, wiseToken)

			defer ioReader.Close()
			wiseResponseData, err = io.ReadAll(ioReader)

			if err != nil {
				return make([]models.Attachment, 0), err
			}
			err = json.Unmarshal(wiseResponseData, &attachmentResponse)
			attachment.Byte = attachmentResponse.Base64Attachment
			attachment.Name = wiseAttachment.Name
			attachment.FileName = wiseAttachment.Name
			attachments = append(attachments, attachment)
		}
		return attachments, err
	}

	var policy models.Policy

	log.Println("Getting attachments for policy saved in Wopta")
	if policy, err = plc.GetPolicy(policyUid, origin); err != nil {
		log.Println("Error when getting policy: " + err.Error())
		return make([]models.Attachment, 0), err
	}

	log.Printf("Found %d attachment(s) for policy %s", len(*policy.Attachments), policy.Uid)

	if policy.Attachments == nil {
		log.Printf("No attachments found for policy %s", policy.Uid)
		return make([]models.Attachment, 0), errors.New("no attachments found for policy " + policy.Uid)
	}

	for _, attachment := range *policy.Attachments {
		if attachment.IsPrivate {
			log.Printf("Attachment %s is private, skipping", attachment.FileName)
			continue
		}
		var responseAttachment models.Attachment
		if len(attachment.Link) == 0 {
			log.Printf("Attachment %s has empty link, skipping", attachment.FileName)
			continue
		}

		gsLink := attachment.Link

		if gsLink == "" {
			log.Printf("empty gsLink")
			continue
		}
		if !strings.Contains(gsLink, "gs://") {
			gsLink = "gs://" + os.Getenv("GOOGLE_STORAGE_BUCKET") + "/" + gsLink
		}

		log.Printf("gsLink: %s", gsLink)

		fileData, err := lib.ReadFileFromGoogleStorage(gsLink)
		if err != nil {
			log.Printf("error reading document from Google Storage: %s", err.Error())
			return nil, err
		}

		responseAttachment.FileName = attachment.FileName
		responseAttachment.ContentType = attachment.ContentType
		responseAttachment.Name = attachment.Name
		responseAttachment.Byte = base64.StdEncoding.EncodeToString(fileData)

		attachments = append(attachments, responseAttachment)
	}

	log.Printf("Sending %d attachment(s)", len(attachments))
	return attachments, err
}

func findNamedMatches(regex *regexp.Regexp, str string) map[string]string {
	match := regex.FindStringSubmatch(str)

	results := map[string]string{}
	for i, name := range match {
		results[regex.SubexpNames()[i]] = name
	}
	return results
}

type WiseAttachmentResponse struct {
	Base64Attachment string `json:"fileAllegato"`
}

type GetPolicyAttachmentsResponse struct {
	Attachments []models.Attachment `json:"attachments"`
}
