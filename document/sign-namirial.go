package document

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"strings"

	"os"
	"path/filepath"
	"time"

	lib "github.com/wopta/goworkspace/lib"
	model "github.com/wopta/goworkspace/models"
)

func SignNamirial(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {

	req := lib.ErrorByte(ioutil.ReadAll(r.Body))
	var data model.Policy
	defer r.Body.Close()
	err := json.Unmarshal([]byte(req), &data)
	lib.CheckError(err)
	//file, _ := os.Open("document/billing.pdf")

	return NamirialOtp(data)

}

type NamirialOtpResponse struct {
	EnvelopeId string `json:"envelopeId"`
	Url        string `json:"url"`
	FileId     string `json:"fileId"`
}

func NamirialOtp(data model.Policy) (string, NamirialOtpResponse, error) {
	var file []byte
	if os.Getenv("env") == "local" {
		file = lib.ErrorByte(ioutil.ReadFile("document/contract.pdf"))

	} else {
		file = lib.GetFromStorage("function-data", data.DocumentName, "")
	}

	var urlstring = os.Getenv("ESIGN_BASEURL") + "v4/sspfile/uploadtemporary"
	//b, _ := ioutil.ReadAll(file)
	SspFileId := <-postData(file, urlstring)
	log.Println(data.Uid+"postData:", SspFileId)
	//prepareEnvelop(SspFileId)
	id := <-sendEnvelop(SspFileId, data)
	log.Println(data.Uid+"sendEnvelop:", id)
	url := <-GetEnvelop(id)
	resp := NamirialOtpResponse{
		EnvelopeId: id,
		Url:        url,
		FileId:     SspFileId,
	}
	return "{}", resp, nil
}
func Autorization() {
	var urlstring = os.Getenv("ESIGN_BASEURL") + "v4/authorization"
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, _ := http.NewRequest(http.MethodGet, urlstring, nil)
	req.Header.Set("apiToken", os.Getenv("ESIGN_TOKEN_API"))
	log.Println("url parse:", req.Header)
	res, err := client.Do(req)
	lib.CheckError(err)

	if res != nil {
		body, err := ioutil.ReadAll(res.Body)
		lib.CheckError(err)
		res.Body.Close()
		log.Println("body:", string(body))
	}
}
func prepareEnvelop(id string) string {
	log.Println("prepare")
	//var b bytes.Buffer
	//fileReader := bytes.NewReader([]byte())
	var urlstring = os.Getenv("ESIGN_BASEURL") + "v4.0/envelope/prepare"
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, _ := http.NewRequest(http.MethodPost, urlstring, strings.NewReader(getPrepare(id)))
	req.Header.Set("apiToken", os.Getenv("ESIGN_TOKEN_API"))

	//header('Content-Length: ' . filesize($pdf));
	log.Println("url parse:", req.Header)
	res, err := client.Do(req)
	lib.CheckError(err)
	var r string
	if res != nil {
		body, err := ioutil.ReadAll(res.Body)
		lib.CheckError(err)
		var result map[string]string
		json.Unmarshal([]byte(body), &result)
		res.Body.Close()
		r = result["SspFileId"]

		log.Println("body:", string(body))
	}

	return r
}
func sendEnvelop(id string, data model.Policy) <-chan string {
	r := make(chan string)

	go func() {
		defer close(r)
		log.Println("Send")
		var urlstring = os.Getenv("ESIGN_BASEURL") + "v4.0/envelope/send"
		client := &http.Client{
			Timeout: time.Second * 10,
		}
		log.Println(getSend(id, data))
		req, _ := http.NewRequest(http.MethodPost, urlstring, strings.NewReader(getSend(id, data)))
		req.Header.Set("apiToken", os.Getenv("ESIGN_TOKEN_API"))
		req.Header.Set("Content-Type", "application/json")
		//header('Content-Length: ' . filesize($pdf));

		res, err := client.Do(req)
		lib.CheckError(err)

		if res != nil {
			body, err := ioutil.ReadAll(res.Body)
			lib.CheckError(err)
			var result map[string]string
			json.Unmarshal([]byte(body), &result)
			res.Body.Close()

			log.Println(data.Uid+" body:", string(body))
			r <- result["EnvelopeId"]

		}
	}()
	return r
}
func GetEnvelop(id string) <-chan string {
	r := make(chan string)

	go func() {
		defer close(r)
		log.Println("Send")
		//var b bytes.Buffer
		//fileReader := bytes.NewReader([]byte())
		var urlstring = os.Getenv("ESIGN_BASEURL") + "/v5/envelope/" + id
		client := &http.Client{
			Timeout: time.Second * 10,
		}
		req, _ := http.NewRequest(http.MethodGet, urlstring, nil)
		req.Header.Set("apiToken", os.Getenv("ESIGN_TOKEN_API"))
		req.Header.Set("Content-Type", "application/json")
		//header('Content-Length: ' . filesize($pdf));

		res, err := client.Do(req)
		lib.CheckError(err)

		if res != nil {
			body, err := ioutil.ReadAll(res.Body)
			lib.CheckError(err)
			var result GetEvelopResponse
			json.Unmarshal([]byte(body), &result)
			res.Body.Close()

			log.Println("body:", string(body))

			r <- result.Bulks[0].Steps[0].WorkstepRedirectionURL

		}
	}()
	return r
}
func postData(data []byte, host string) <-chan string {
	r := make(chan string)

	go func() {
		defer close(r)
		var b bytes.Buffer
		w := multipart.NewWriter(&b)
		// Add the field
		fw, err := w.CreateFormFile("file", filepath.Base("contract.pdf"))
		lib.CheckError(err)
		fw.Write((data)[:])
		w.Close()
		log.Println("Post")
		req, err := http.NewRequest("POST", host, &b)
		lib.CheckError(err)
		req.Header.Set("apiToken", os.Getenv("ESIGN_TOKEN_API"))
		req.Header.Set("Content-Type", w.FormDataContentType())
		client := &http.Client{
			Timeout: time.Second * 10,
		}
		res, err := client.Do(req)
		var result map[string]string
		resByte, err := ioutil.ReadAll(res.Body)
		lib.CheckError(err)
		json.Unmarshal(resByte, &result)
		res.Body.Close()
		log.Println("Post 2")
		r <- result["SspFileId"]
		fmt.Println(res.StatusCode)
	}()

	return r
}
func GetFile(id string, uid string) chan string {
	r := make(chan string)

	go func() {
		defer close(r)
		var urlstring = os.Getenv("ESIGN_BASEURL") + "v6/file/" + id
		client := &http.Client{
			Timeout: time.Second * 10,
		}
		req, _ := http.NewRequest(http.MethodGet, urlstring, nil)
		req.Header.Set("apiToken", os.Getenv("ESIGN_TOKEN_API"))
		log.Println("url parse:", req.Header)
		res, err := client.Do(req)
		lib.CheckError(err)
		if res != nil {
			body, err := ioutil.ReadAll(res.Body)
			lib.CheckError(err)
			log.Println(body)
			res.Body.Close()
			lib.PutToFireStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "/document/contracts/"+uid, body)
			r <- "upload done"

		}

	}()
	return r
}
func getPrepare(id string) string {
	return `{
	"SspFileIds": [
	  " ` + id + `"
	],
	"AdHocWorkstepConfiguration": {
	  "WorkstepLabel": "string",
	  "SmallTextZoomFactorPercent": 0,
	  "WorkstepTimeToLiveInMinutes": 0,
	  "FinishAction": {
		"ClientActions": [
		  {
			"RemoveDocumentFromRecentDocumentList": true,
			"CallClientActionOnlyAfterSuccessfulSync": true,
			"ClientName": "string",
			"CloseApp": true,
			"Action": "string"
		  }
		]
	  },
	  "NoSequenceEnforced": true,
	  "SigTemplate": {
		"Size": {
		  "Height": 0,
		  "Width": 0
		},
		"AllowedSignatureTypes": [
		  
		]
	  },
	  "ParseFormFields": {
		"MapRequiredFieldsToRequiredTask": true,
		"FormsGrouping": "PerPage",
		"ReturnSimplifiedConfig": true,
		"AddKeepExistingValueFlag": true,
		"ParseFormField": true
	  },
	  "AdhocPolicies": {
		"AllowModificationsAfterSignature": true
	  },
	  "ViewerPreferences": {
		"ShowPageNavigationBar": true,
		"ShowThumbnails": true,
		"SkipFinishConfirmDialog": true,
		"SkipDocumentDialog": true,
		"ShowImagesInFullWidth": true,
		"DisableGeolocation": true,
		"ShowDocumentDownloadDialogAfterAutomaticFinish": true,
		"AttachmentsMaxFileSize": 0,
		"SkipPreviewImageOnDisposableCertificate": true,
		"LoadCustomJs": true,
		"AllowCustomButtons": true,
		"GuidingBehavior": "GuideOnlyRequiredTasks",
		"FormFieldsGuidingBehavior": "AllowSubmitAlways",
		"ShowVersionNumber": true,
		"EnableWarningPopupOnLeave": true,
		"WarningPopupDisplayAfter": "FillOrSignField",
		"FinishWorkstepOnOpen": true,
		"AutoFinishAfterRequiredTasksDone": true,
		"GuidingBehaviorOnFinishedTask": "NoMove",
		"SkipThankYouDialog": true,
		"NativeAppsUrlScheme": "string",
		"DocumentViewingMode": "EndlessPaperAllDocuments",
		"ThumbnailMode": "ShowAllPages",
		"ShowTopBar": true,
		"DisplayRejectButtonInTopBar": true,
		"MultipleSignatureTypesAndBatchSigningSettings": {
		  "IsUseBatchSigningCheckedByDefault": true,
		  "IsRememberSignatureTypeCheckedByDefault": true,
		  "IsRememberBatchSigningDecisionCheckedByDefault": true,
		  "SkipMultipleSignatureTypesAndBatchSigningDialogIfBatchSigningPossible": true
		},
		"VisibleAreaOptions": {
		  "AllowedDomain": "string",
		  "Enabled": true
		},
		"ShowStartGuidingHint": true,
		"ShowStatusBar": true,
		"ShowZoomButtons": true,
		"ShowNoGeolocationWarning": true,
		"AutoStartGuiding": true,
		"ShowPageGap": true,
		"ShowPageNavigationButtons": true,
		"ShowFinishPossibleHint": true,
		"SkipRejectConfirmDialog": true,
		"BatchSigningType": "Basic",
		"BatchSigningDisableNextButtonUntilDialogScrolledToBottom": true
	  },
	  "SignatureConfigurations": [
		{
		  "SpcId": "string",
		  "PdfSignatureProperties": {
			"PdfAConformant": true,
			"PAdESPart4Compliant": true,
			"IncludeSigningCertificateChain": true,
			"SigningCertificateRevocationInformationIncludeMode": "DoNotInclude",
			"SignatureTimestampData": {
			  "Uri": "string",
			  "Username": "string",
			  "Password": "string",
			  "SignatureHashAlgorithm": "Sha1",
			  "AuthenticationCertifiateDescriptor": {
				"Identifier": "string",
				"Type": "string"
			  }
			},
			"EnableEutlVerification": true,
			"EnableValidateSigningCertificateName": true,
			"SigningCertificateNameRegex": "string"
		  },
		  "PdfSignatureCryptographicData": {
			"SignatureHashAlgorithm": "Sha1",
			"SigningCertificateDescriptor": {
			  "Identifier": "string",
			  "Type": "Sha1Thumbprint",
			  "Csp": "Default"
			}
		  },
		  "CertificateFilter": {
			"KeyUsages": [
			  "string"
			],
			"ThumbPrints": [
			  "string"
			],
			"RootThumbPrints": [
			  "string"
			]
		  }
		}
	  ],
	  "SigStringParsingConfiguration": {
		"SigStringsForParsings": [
		  {
			"StartPattern": "string",
			"EndPattern": "string",
			"ClearSigString": true,
			"SearchEntireWordOnly": true
		  }
		]
	  },
	  "GeneralPolicies": {
		"AllowSaveDocument": true,
		"AllowSaveAuditTrail": true,
		"AllowRotatingPages": true,
		"AllowAppendFileToWorkstep": true,
		"AllowAppendTaskToWorkstep": true,
		"AllowEmailDocument": true,
		"AllowPrintDocument": true,
		"AllowFinishWorkstep": true,
		"AllowRejectWorkstep": true,
		"AllowRejectWorkstepDelegation": true,
		"AllowUndoLastAction": true,
		"AllowColorizePdfForms": true,
		"AllowAdhocPdfAttachments": true,
		"AllowAdhocSignatures": true,
		"AllowAdhocStampings": true,
		"AllowAdhocFreeHandAnnotations": true,
		"AllowAdhocTypewriterAnnotations": true,
		"AllowAdhocPictureAnnotations": true,
		"AllowAdhocPdfPageAppending": true,
		"AllowReloadOfFinishedWorkstep": true
	  },
	  "FinalizeActions": {
		
	  },
	  "TransactionCodeConfigurations": [
		{
		  "Id": "string",
		  "HashAlgorithmIdentifier": "Sha1",
		  "Texts": [
			{
			  "Language": "string",
			  "Value": "string"
			}
		  ]
		}
	  ]
	},
	"PrepareSendEnvelopeStepsDescriptor": {
	  "ClearFieldMarkupString": true
	}
  }`
}

func getSend(id string, data model.Policy) string {

	return `{
  
		"SspFileIds": [
		  " ` + id + `"
		],
		"SendEnvelopeDescription":{
			"Name": "` + id + `.pdf",
			"EmailSubject": "Please sign the enclosed envelope",
			"EmailBody": "Dear #RecipientFirstName# #RecipientLastName#\n\n#PersonalMessage#\n\nPlease sign the envelope #EnvelopeName#\n\nEnvelope will expire at #ExpirationDate#",
			"DisplayedEmailSender": "",
			"EnableReminders": true,
			"FirstReminderDayAmount": 5,
			"RecurrentReminderDayAmount": 3,
			"BeforeExpirationDayAmount": 3,
			"ExpirationInSecondsAfterSending": 2419200,
			"CallbackUrl": "",
			"StatusUpdateCallbackUrl": "https://europe-west1-positive-apex-350507.cloudfunctions.net/callback/v1/sign?envelope=##EnvelopeId##&action=##Action##&uid=` + data.Uid + `&token=` + os.Getenv("WOPTA_TOKEN_API") + `",
			"LockFormFieldsAtEnvelopeFinish": false,
			"Steps": [
			  {
				"OrderIndex": 1,
				"Recipients": [
				  {
					"Email": "` + data.Contractor.Mail + `",
					"FirstName": "` + data.Contractor.Name + `",
					"LastName": "` + data.Contractor.Surname + `",
					"LanguageCode": "it",
					"EmailBodyExtra": "",
					"DisableEmail": false,
					"AddAndroidAppLink": false,
					"AddIosAppLink": false,
					"AddWindowsAppLink": false,
					"AllowDelegation": true,
					"AllowAccessFinishedWorkstep": false,
					"SkipExternalDataValidation": false,
					"AuthenticationMethods": [{
						"Method": "Sms",
						"Parameter": "` + data.Contractor.Phone + `"
		   }],
					"IdentificationMethods": [],
					"OtpData": {
					  "PhoneMobile": "` + data.Contractor.Phone + `"
					}
				  }
				],
				"EmailBodyExtra": "",
				"RecipientType": "Signer",
				"WorkstepConfiguration": {
				  "WorkstepLabel": "` + data.Contractor.Name + " " + data.Contractor.Surname + " " + data.Name + `",
				  "SmallTextZoomFactorPercent": 100,
				  "FinishAction": {
					"ServerActions": [],
					"ClientActions": []
				  },
				  "ReceiverInformation": {
					"UserInformation": {
					  "FirstName": "` + data.Contractor.Name + `",
					  "LastName": "` + data.Contractor.Surname + `",
					  "EMail": "` + data.Contractor.Mail + `"
					},
					"TransactionCodePushPluginData": []
				  },
				  "SenderInformation": {
					"UserInformation": {
					  "FirstName": "Wopta",
					  "LastName": "Assicurazzioni",
					  "EMail": "info@wopta.it"
					}
				  },
				  "TransactionCodeConfigurations": [],
				  "SignatureConfigurations": [],
				  "ViewerPreferences": {
					"FinishWorkstepOnOpen": false,
					"VisibleAreaOptions": {
					  "AllowedDomain": "",
					  "Enabled": false
					}
				  },
				  "ResourceUris": {
					"DelegationUri": ""
				  },
				  "AuditingToolsConfiguration": {
					"WriteAuditTrail": true
				  },
				  "Policy": {
					"GeneralPolicies": {
					  "AllowSaveDocument": true,
					  "AllowSaveAuditTrail": true,
					  "AllowRotatingPages": false,
					  "AllowAppendFileToWorkstep": false,
					  "AllowAppendTaskToWorkstep": false,
					  "AllowEmailDocument": true,
					  "AllowPrintDocument": true,
					  "AllowFinishWorkstep": true,
					  "AllowRejectWorkstep": true,
					  "AllowRejectWorkstepDelegation": true,
					  "AllowUndoLastAction": true,
					  "AllowColorizePdfForms": false,
					  "AllowAdhocPdfAttachments": false,
					  "AllowAdhocSignatures": false,
					  "AllowAdhocStampings": false,
					  "AllowAdhocFreeHandAnnotations": false,
					  "AllowAdhocTypewriterAnnotations": false,
					  "AllowAdhocPictureAnnotations": false,
					  "AllowAdhocPdfPageAppending": false,
					  "AllowReloadOfFinishedWorkstep": true
					},
					"WorkstepTasks": {
					  "PictureAnnotationMinResolution": 0,
					  "PictureAnnotationMaxResolution": 0,
					  "PictureAnnotationColorDepth": "Color16M",
					  "SequenceMode": "NoSequenceEnforced",
					  "PositionUnits": "PdfUnits",
					  "ReferenceCorner": "Lower_Left",
					  "Tasks": [
						{
						  "Texts": [
							{
							  "Language": "it",
							  "Value": "Signature Disclosure Text"
							},
							{
							  "Language": "*",
							  "Value": "Signature Disclosure Text"
							}
						  ],
						  "Headings": [
							{
							  "Language": "it",
							  "Value": "Signature Disclosure Subject"
							},
							{
							  "Language": "*",
							  "Value": "Signature Disclosure Subject"
							}
						  ],
						  "IsRequired": false,
						  "Id": "ra",
						  "DisplayName": "ra",
						  "DocRefNumber": 1,
						  "DiscriminatorType": "Agreements"
						},
						{
						  "PositionPage": 2,
						  "Position": {
							"PositionX": 252.0,
							"PositionY": 186.0
						  },
						  "Size": {
							"Height": 80.0,
							"Width": 190.0
						  },
						  "AdditionalParameters": [
							{
							  "Key": "enabled",
							  "Value": "1"
							},
							{
							  "Key": "completed",
							  "Value": "0"
							},
							{
							  "Key": "req",
							  "Value": "1"
							},
							{
							  "Key": "isPhoneNumberRequired",
							  "Value": "0"
							},
							{
							  "Key": "trValidityInSeconds",
							  "Value": "60"
							},
							{
							  "Key": "fd",
							  "Value": ""
							},
							{
							  "Key": "fd_dateformat",
							  "Value": "dd-MM-yyyy HH:mm:ss"
							},
							{
							  "Key": "fd_timezone",
							  "Value": "datetimeutc"
							}
						  ],
						  "AllowedSignatureTypes": [
							{
							  "TrModType": "TransactionCodeSenderPlugin",
							  "TrValidityInSeconds": 300,
							  "TrConfId": "otpSignatureSmsText",
							  "IsPhoneNumberRequired": true,
							  "Ly": "simpleTransactionCodeSms",
							  "Id": "c787919a-b2fd-4849-8f97-98dee281da30",
							  "DiscriminatorType": "SigTypeTransactionCode",
							  "Preferred": false,
							  "StampImprintConfiguration": {
								"DisplayExtraInformation": true,
								"DisplayEmail": true,
								"DisplayIp": true,
								"DisplayName": true,
								"DisplaySignatureDate": true,
								"FontFamily": "Times New Roman",
								"FontSize": 11.0,
								"OverrideLegacyStampImprint": false,
								"DisplayTransactionId": true,
								"DisplayTransaktionToken": true,
								"DisplayPhoneNumber": true
							  },
							  "SignaturePluginConfigurationId": "ltaLevelId"
							}
						  ],
						  "UseTimestamp": false,
						  "IsRequired": true,
						  "Id": "1#XyzmoDuplicateIdSeperator#Signature_e7ca3f6a-33fa-cdba-d696-1377fcad51c9",
						  "DisplayName": "",
						  "DocRefNumber": 1,
						  "DiscriminatorType": "Signature"
						}
					  ]
					},
					"FinalizeActions": {
					  "FinalizeActionList": [
						{
						  "DocRefNumbers": "*",
						  "SpcId": "ltaLevelId",
						  "DiscriminatorType": "Timestamp"
						}
					  ]
					}
				  },
				  "Navigation": {
					"HyperLinks": [],
					"Links": [],
					"LinkTargets": []
				  }
				},
				"DocumentOptions": [
				  {
					"DocumentReference": "1",
					"IsHidden": false
				  }
				],
				"UseDefaultAgreements": true
			  },
			  
			],
			"AddFormFields": {
			  "Forms": {}
			},
			"OverrideFormFieldValues": {
			  "Forms": {}
			},
			"AttachSignedDocumentsToEnvelopeLog": false
		  }
	  }`
}
func getSendTemplate(id string) string {
	return `{
		"TemplateId": "string",
		"EnvelopeOverrideOptions": {
		  "Recipients": [
			{
			  "RecipientId": "string",
			  "OrderIndex": 0,
			  "Email": "string",
			  "Recipient": {
				"Email": "string",
				"FirstName": "string",
				"LastName": "string",
				"LanguageCode": "string",
				"EmailBodyExtra": "string",
				"DisableEmail": true,
				"AddAndroidAppLink": true,
				"AddIosAppLink": true,
				"AddWindowsAppLink": true,
				"AllowDelegation": true,
				"AllowAccessFinishedWorkstep": true,
				"SkipExternalDataValidation": true,
				"AuthenticationMethods": [
				  {
					"Method": "Pin",
					"Parameter": "string",
					"Filters": [
					  {
						"CompareOperation": "Equals",
						"FilterId": "string",
						"FilterValue": "string"
					  }
					]
				  }
				],
				"IdentificationMethods": [
				  {
					"Method": "OAuth",
					"Parameter": "string",
					"Filters": [
					  {
						"CompareOperation": "Equals",
						"FilterId": "string",
						"FilterValue": "string"
					  }
					]
				  }
				],
				"DisposableCertificateData": {
				  "CountryResidence": "string",
				  "DocumentIssuingCountry": "string",
				  "IdentificationCountry": "string",
				  "IdentificationType": "NONE",
				  "PhoneMobile": "string",
				  "DocumentType": "CI",
				  "DocumentIssuedBy": "string",
				  "DocumentIssuedOn": "2022-10-27T14:15:24.573Z",
				  "DocumentExpiryDate": "2022-10-27T14:15:24.573Z",
				  "SerialNumber": "string",
				  "DocumentNumber": "string",
				  "OverrideHolderInCaseOfMismatch": true
				},
				"SwissComCertificateData": {
				  "PhoneNumber": "string",
				  "Parameters": [
					{
					  "Key": "string",
					  "Value": "string"
					}
				  ]
				},
				"RemoteCertificateData": {
				  "UserId": "string",
				  "DeviceId": "string"
				},
				"OtpData": {
				  "PhoneMobile": "string"
				},
				"Pkcs7SignerData": {
				  "AllowedPkcs7SignatureTypes": [
					"LocalCertificate"
				  ]
				}
			  }
			}
		  ],
		  "AddFormFields": {
			"Forms": {}
		  },
		  "OverrideFormFieldValues": {
			"Forms": {}
		  },
		  "Name": "string",
		  "EmailSubject": "string",
		  "EmailBody": "string",
		  "EnableReminders": true,
		  "FirstReminderDayAmount": 0,
		  "RecurrentReminderDayAmount": 0,
		  "BeforeExpirationDayAmount": 0,
		  "DaysUntilExpire": 0,
		  "ExpirationDate": "2022-10-27T14:15:24.573Z",
		  "ExpirationInSecondsAfterSending": 0,
		  "CallbackUrl": "string",
		  "StatusUpdateCallbackUrl": "string",
		  "WorkstepEventCallback": {
			"Url": "string",
			"Blacklist": [
			  "string"
			],
			"WhiteList": [
			  "string"
			]
		  },
		  "MetaDataXml": "string"
		}
	  }`
}

func UnmarshalWelcome(data []byte) (GetEvelopResponse, error) {
	var r GetEvelopResponse
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *GetEvelopResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type GetEvelopResponse struct {
	Status                          string `json:"Status"`
	SendDate                        string `json:"SendDate"`
	ExpirationDate                  string `json:"ExpirationDate"`
	ValidityFromCreationInDays      int64  `json:"ValidityFromCreationInDays"`
	ExpirationInSecondsAfterSending int64  `json:"ExpirationInSecondsAfterSending"`
	Bulks                           []Bulk `json:"Bulks"`

	ID   string `json:"Id"`
	Bulk string `json:"Bulk"`

	LockFormFieldsAtEnvelopeFinish bool `json:"LockFormFieldsAtEnvelopeFinish"`
}

type Bulk struct {
	Status            string        `json:"Status"`
	Email             string        `json:"Email"`
	LogDocumentID     string        `json:"LogDocumentId"`
	FinishedDocuments []interface{} `json:"FinishedDocuments"`
	Steps             []Step        `json:"Steps"`
}

type Step struct {
	ID                          string `json:"Id"`
	FirstName                   string `json:"FirstName"`
	LastName                    string `json:"LastName"`
	OrderIndex                  int64  `json:"OrderIndex"`
	Email                       string `json:"Email"`
	LanguageCode                string `json:"LanguageCode"`
	Status                      string `json:"Status"`
	StatusReason                string `json:"StatusReason"`
	RecipientType               string `json:"RecipientType"`
	WorkstepRedirectionURL      string `json:"WorkstepRedirectionUrl"`
	AllowAccessFinishedWorkstep bool   `json:"AllowAccessFinishedWorkstep"`

	IsParallel bool `json:"IsParallel"`
}
