package policy

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/go-chi/chi/v5"

	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/wiseproxy"
)

func getPolicyAttachmentsFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		response GetPolicyAttachmentsResponse
	)

	log.AddPrefix("GetPolicyAttachmentsFx")
	defer log.PopPrefix()

	log.Println("Handler start -----------------------------------------------")

	attachments, err := GetPolicyAttachments(chi.URLParam(r, "uid"))
	if err != nil {
		log.Println("GetPolicyAttachments Error: " + err.Error())
		return "{}", nil, nil
	}

	response.Attachments = attachments
	res, err := json.Marshal(response)
	if err != nil {
		log.Println("AttachmentsMarshal Error: " + err.Error())
	}

	log.Println("Handler end -------------------------------------------------")

	return string(res), nil, nil
}

func GetPolicyAttachments(policyUid string) ([]models.Attachment, error) {
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
		request := fmt.Appendf(nil, `{"idPolizza": "%s", "cdLingua": "it"}`, wiseAttachmentId)
		ioReader := wiseproxy.WiseProxyObj("WebApiProduct/Api/GetPolizzaCompleta", request, http.MethodPost)

		defer ioReader.Close()
		wiseResponseData, err = io.ReadAll(ioReader)

		if err != nil {
			return make([]models.Attachment, 0), err
		}
		err = json.Unmarshal(wiseResponseData, &wisePolicy)

		for _, wiseAttachment := range wisePolicy.Policy.Attachments {
			var attachment models.Attachment

			request := fmt.Appendf(nil, `{"txRifAllegato": "%s", "cdLingua": "it"}`, wiseAttachment.Id)
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
	if policy, err = GetPolicy(policyUid); err != nil {
		log.Println("Error when getting policy: " + err.Error())
		return make([]models.Attachment, 0), err
	}

	expr, err := regexp.Compile("gs://(?P<bucketName>(?:[^/])*)/(?P<fileName>((?:[^/]*/)*)(.*))")
	log.Printf("Found %d attachment(s) for policy %s", len(*policy.Attachments), policy.Uid)
	for _, attachment := range *policy.Attachments {
		var responseAttachment models.Attachment
		if len(attachment.Link) == 0 {
			log.Printf("Attachment %s has empty link, skipping", attachment.FileName)
			continue
		}
		matches := findNamedMatches(expr, attachment.Link)
		log.Printf("Found %s with bucketName=%s and fileName=%s", attachment.FileName, matches["bucketName"], matches["fileName"])
		fileData, _ := lib.GetFromStorageErr(matches["bucketName"], matches["fileName"], "")

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
