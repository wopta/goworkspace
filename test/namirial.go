package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/wopta/goworkspace/document"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	"github.com/wopta/goworkspace/policy"
	"github.com/wopta/goworkspace/product"
)

type fileDesc struct {
	bite []byte
	desc string
	FileId   string
	DocumentNumber  int 
}
type requestEnvelop struct{
	Uids []struct {
		Uid string
	}
}
func HandlerEnvelop(w http.ResponseWriter, r *http.Request) (string, any, error) {
	log.SetPrefix("[CreateEnvelop]")
	defer log.SetPrefix("")
	
	req,err:=getBodyRequest[requestEnvelop](r)
	if err!=nil{
		return "",nil,err
	}

	files := make([]fileDesc, 0)
	var name string
	var policyModel models.Policy
	var networkModel *models.NetworkNode
	var productModel *models.Product
	var warrant *models.Warrant
	for i :=range req.Uids{

		log.Println("Creation uid policy ",req.Uids[i].Uid)
		policyModel,err=policy.GetPolicy(req.Uids[i].Uid,"")
		networkModel=network.GetNetworkNodeByUid(policyModel.ProducerUid)
		if networkModel!=nil{
			warrant=networkModel.GetWarrant()
		}
		productModel=product.GetProductV2(policyModel.Name,policyModel.ProductVersion,policyModel.Channel,networkModel,warrant)
		docu:=document.Proposal("",&policyModel,networkModel,productModel)

		log.Println("Policy company ",policyModel.Company)
		if err!=nil{
			return "",nil,err
		}

		log.Println("Document status ",docu.LinkGcs)

		SspFileId, err := uploadFile([]byte(docu.Bytes), name)
		if err != nil {
			return "", "", err
		}
		log.Printf("File id %v", SspFileId)
		files = append(files,
		fileDesc{
			bite: []byte(docu.Bytes),
			desc: fmt.Sprint("constratto", i),
			FileId:   SspFileId,
			DocumentNumber: i,
		})
	}

	res,err:=sendEnvelop(files)
	if err!=nil{
		return "",nil,err
	}
	log.Println("Envelop id",res)

	return res,nil,nil
}

func uploadFile(bites []byte, name string) (string, error) {
	log.SetPrefix("[UploadFile")
	defer log.SetPrefix("")

	var urlstring = os.Getenv("ESIGN_BASEURL") + "v6/file/upload"
	log.Println("Url ", urlstring)

	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)
	defer writer.Close()

	fw, err := writer.CreateFormFile("file", name)
	if _, e := fw.Write(bites); e != nil {
		panic(e)
	}

	req, err := http.NewRequest("POST", urlstring, &buffer)
	lib.CheckError(err)
	req.Header.Set("apiToken", os.Getenv("ESIGN_TOKEN_API"))
	req.Header.Set("Authorization", "Bearer"+os.Getenv("ESIGN_TOKEN_API"))
	req.Header.Set("Content-Type", writer.FormDataContentType())

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

	filesJson,err:=json.Marshal(files)
	request:=getRequestTOSend(string(filesJson))
	log.Printf("Body for request:",request)

	req, err := http.NewRequest("POST", urlstring, bytes.NewReader([]byte(request)))
	lib.CheckError(err)
	req.Header.Set("apiToken", os.Getenv("ESIGN_TOKEN_API"))
	req.Header.Set("Content-Type", "application/json")

	log.Println("Call ...")
	res, err := lib.RetryDo(req, 5, 30)
	log.Println("Result ...", res.Status)
	if res.StatusCode==http.StatusOK{
		body,_:=getBodyResponse[struct {EvelopeId string}](res)
		return body.EvelopeId,nil
	}
	return "",fmt.Errorf("errore invio",res.Status)
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
							"LanguageCode": "EN"
						}
					},
					"Elements": {
						"Signatures": [{
							"ElementId": "sample sig click2sign",
							"Required": true,
							"DocumentNumber": 1,
							"DisplayName": "Sign here",
							"AllowedSignatureTypes": {
								"ClickToSign": {
								}
							},
							"FieldDefinition": {
								"Position": {
									"PageNumber": 1,
									"X": 100,
									"Y": 200
								},
								"Size": {
									"Width": 100,
									"Height": 70
								}
							}
						}
						]
					},
					"SigningGroup": "firstSigner"
				}
			}
		}, {
			"Action": {
				"SendCopy": {
					"RecipientConfiguration": {
						"ContactInformation": {
							"Email": "john.doe@sample.com",
							"GivenName": "John",
							"Surname": "Doe",
							"LanguageCode": "EN"
						}
					}
				}
			}
		}
		]
	}`, filesJson)
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
