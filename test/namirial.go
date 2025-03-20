package test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/wopta/goworkspace/document"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	"github.com/wopta/goworkspace/policy"
	"github.com/wopta/goworkspace/product"
)

type fileDesc struct {
	bite           []byte
	desc           string
	FileId         string
	DocumentNumber int
}
type requestEnvelop struct {
	Uids []struct {
		Uid string
	}
}

type responseEnvelop struct {
	EnvelopeId string
}

func HandlerEnvelop(w http.ResponseWriter, r *http.Request) (string, any, error) {
	log.SetPrefix("[CreateEnvelop]")
	defer log.SetPrefix("")

	req, err := getBodyRequest[requestEnvelop](r)
	if err != nil {
		return "", nil, err
	}

	files := make([]fileDesc, 0)
	var policyModel models.Policy
	var networkModel *models.NetworkNode
	var productModel *models.Product
	var warrant *models.Warrant
	for i := range req.Uids {

		log.Println("Creation uid policy ", req.Uids[i].Uid)
		policyModel, err = policy.GetPolicy(req.Uids[i].Uid, "")
		networkModel = network.GetNetworkNodeByUid(policyModel.ProducerUid)
		if networkModel != nil {
			warrant = networkModel.GetWarrant()
		}
		productModel = product.GetProductV2(policyModel.Name, policyModel.ProductVersion, policyModel.Channel, networkModel, warrant)
		docu := document.Proposal("", &policyModel, networkModel, productModel)
		log.Println("Policy company ", policyModel.Company)
		if err != nil {
			return "", nil, err
		}

		if docu == nil {
			return "Error creating file", nil, err
		}
		byteDoc, err := base64.StdEncoding.DecodeString(docu.Bytes)
		if err != nil {
			return "", "", err
		}
		if len(byteDoc) == 0 {
			return "Error creating file, no bytes", nil, err
		}
		SspFileId, err := uploadFile(byteDoc, "Contranct "+fmt.Sprint(i)+".pdf")
		if err != nil {
			return "", "", err
		}
		if SspFileId == "" {
			return "", nil, fmt.Errorf("Error uploading file, no file id found")
		}
		log.Printf("File id %v", SspFileId)
		files = append(files,
		fileDesc{
			bite:           []byte(byteDoc),
			desc:           fmt.Sprint("constratto", i),
			FileId:         SspFileId,
			DocumentNumber: i + 1,
		})
	}

	resJson, err := sendEnvelop(files)
	if err != nil {
		return "", nil, err
	}
	log.Println("Envelop id", resJson)
	var resEnvelop responseEnvelop
	json.Unmarshal([]byte(resJson), &resEnvelop)

	resJson, err = getEnveloper(resEnvelop.EnvelopeId)
	if err != nil {
		return "", nil, err
	}
	return resJson, nil, nil
}

func uploadFile(bites []byte, name string) (string, error) {
	log.SetPrefix("[UploadFile")
	defer log.SetPrefix("")

	var urlstring = os.Getenv("ESIGN_BASEURL") + "v6/file/upload"
	log.Println("Url ", urlstring)

	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	fw, err := writer.CreateFormFile("file", name+".pdf")
	if _, e := fw.Write((bites)[:]); e != nil {
		panic(e)
	}

	writer.Close()

	req, err := http.NewRequest("POST", urlstring, &buffer)
	lib.CheckError(err)
	req.Header.Set("apiToken", os.Getenv("ESIGN_TOKEN_API"))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	//req.Header.Set("Authorization", "Bearer "+os.Getenv("ESIGN_TOKEN_API"))

	log.Println("Call ...")
	res, err := lib.RetryDo(req, 5, 30)
	log.Println("Result ...", res.Status)

	if res.StatusCode == http.StatusOK {
		body, err := getBodyResponse[struct{ FileId string }](res)
		if err != nil {
			return "", err
		}
		return body.FileId, nil
	}
	return "", nil
}

func sendEnvelop(files []fileDesc) (string, error) {
	log.SetPrefix("[SendEnvelop]")
	defer log.SetPrefix("")

	var urlstring = os.Getenv("ESIGN_BASEURL") + "v6/envelope/send"
	log.Println("Url ", urlstring)

	filesJson, err := json.Marshal(files)
	request := getRequestTOSend(string(filesJson))
	log.Println("Body for request:")
	log.Println(request)

	req, err := http.NewRequest("POST", urlstring, strings.NewReader(request))
	lib.CheckError(err)
	req.Header.Set("apiToken", os.Getenv("ESIGN_TOKEN_API"))
	req.Header.Set("Content-Type", "application/json")

	log.Println("Call ...")
	res, err := lib.RetryDo(req, 5, 30)
	log.Println("Result ...", res.Status)
	body := lib.ErrorByte(io.ReadAll(res.Body))
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		return string(body), nil
	}
	return "", fmt.Errorf(res.Status, string(body))
}
func getRequestTOSend(filesJson string) string {
	return fmt.Sprintf(`{
		"Documents": %v,
		"Name": "Test",
		"Activities": [{
			"Action": {
				"Sign": {
					"RecipientConfiguration": {
						"ContactInformation": {
							"Email": "jane.doe@sample.com",
							"GivenName": "Jane",
							"Surname": "Doe",
							"LanguageCode": "it"
						}
					},
					"Elements": {
						"Signatures": [{
							"elementid": "sample sig clicksign",
							"required": true,
							"documentnumber": 1,
							"displayname": "sign here 1",
							"allowedsignaturetypes": {
								"clicktosign": {
								}
							},
							"fielddefinition": {
								"position": {
									"pagenumber": 1,
									"x": 100,
									"y": 200
								},
								"size": {
									"width": 100,
									"height": 70
								}
							}
						}
						]
					},
					"SigningGroup": "firstSigner"
				},
				"Sign": {
					"RecipientConfiguration": {
						"ContactInformation": {
							"Email": "jane.doe@sample.com",
							"GivenName": "Jane",
							"Surname": "Doe",
							"LanguageCode": "it"
						}
					},
					"Elements": {
						"Signatures": [{
							"elementid": "sample sig click2sign",
							"required": true,
							"documentnumber": 2,
							"displayname": "sign here 2",
							"allowedsignaturetypes": {
								"clicktosign": {
								}
							},
							"fielddefinition": {
								"position": {
									"pagenumber": 1,
									"x": 100,
									"y": 200
								},
								"size": {
									"width": 100,
									"height": 70
								}
							}
						}
						]
					},
					"SigningGroup": "secondSigner"
				}
			}
		}
		]
	}`, filesJson)
}
func getEnveloper(envelopeId string) (string, error) {
	log.SetPrefix("[GetEnvelop]")
	defer log.SetPrefix("")

	log.Println("Envelope to get: ", envelopeId)
	var urlstring = fmt.Sprint(os.Getenv("ESIGN_BASEURL")+"v6/envelope/"+envelopeId, "/viewerlinks")

	log.Println("Url ", urlstring)

	req, err := http.NewRequest("GET", urlstring, nil)
	lib.CheckError(err)

	req.Header.Set("apiToken", os.Getenv("ESIGN_TOKEN_API"))
	req.Header.Set("Content-Type", "none")

	log.Println("Call ...")
	res, err := lib.RetryDo(req, 5, 30)
	log.Println("Result ...", res.Status)
	body := lib.ErrorByte(io.ReadAll(res.Body))
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		return string(body), nil
	}
	return "", fmt.Errorf(res.Status, string(body))
}
func getBodyRequest[T any](r *http.Request) (T, error) {
	var (
		req T
	)
	log.SetPrefix("[Unmarshal] ")
	defer log.SetPrefix("")

	body := lib.ErrorByte(io.ReadAll(r.Body))

	defer r.Body.Close()
	if err := json.Unmarshal(body, &req); err != nil {
		return *new(T), err
	}
	return req, nil
}

func getBodyResponse[T any](r *http.Response) (T, error) {
	var (
		req T
	)
	log.SetPrefix("[Unmarshal] ")
	defer log.SetPrefix("")

	body := lib.ErrorByte(io.ReadAll(r.Body))

	defer r.Body.Close()
	if err := json.Unmarshal(body, &req); err != nil {
		return *new(T), err
	}
	return req, nil
}
