package main

// GOOGLE_CLOUD_PROJECT is a user-set environment variable.
/*{
	"BaseUrl":"https://api-devexternal.munichre.com/flowin/dev/api/V1",
	"ApimKey":"59c92bc0095d4b8c803656a207150c32",
	"TokenEndPoint":"https://login.microsoftonline.com/9f2c9c2d-da50-4f33-8dfb-a780f38b50dd/oauth2/v2.0/token",
	"Scope":"46e8daaf-f894-464a-942a-e06852ed4526/.default",
	"ClientId":"194d46f8-0779-4e17-a96d-62c7bdd81901",
	"ClientSecret":"nrDLtgtLiVhvaChj1sU7JiCBUZbztRXw2ROMBYxZ",
	"GrantType":"client_credentials",
	"UWRole":"Agent",
	"SubProductId_PMIW":"35"
}*/
import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	enr "github.com/wopta/goworkspace/enrich-vat"
	q "github.com/wopta/goworkspace/quote-allrisk"
	rules "github.com/wopta/goworkspace/rules"
)

func main() {
	r := mux.NewRouter().StrictSlash(true)
	os.Setenv("munichreBaseUrl", "https://api-devexternal.munichre.com/flowin/dev/api/V1")
	os.Setenv("munichreTokenEndPoint", "https://login.microsoftonline.com/9f2c9c2d-da50-4f33-8dfb-a780f38b50dd/oauth2/v2.0/token")
	os.Setenv("munichreScope", "46e8daaf-f894-464a-942a-e06852ed4526/.default")
	os.Setenv("munichreClientId", "194d46f8-0779-4e17-a96d-62c7bdd81901")
	os.Setenv("munichreClientSecret", "nrDLtgtLiVhvaChj1sU7JiCBUZbztRXw2ROMBYxZ")
	os.Setenv("munichreSubscriptionKey", "59c92bc0095d4b8c803656a207150c32")
	os.Setenv("munichreSubscriptionHeader", "Ocp-Apim-Subscription-Key")
	os.Setenv("SA_KEY", `{
		"type": "service_account",
		"project_id": "positive-apex-350507",
		"private_key_id": "3df210b6ea35a958f1556992720cd7b176d75e99",
		"private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQCRQqz18GQJ+C/Z\n6CsikwUtGu10Hrxgg2BTEV9nQrZ2opgF8u7HjHnm3rQ94wlxP4vC2lPfnN/gQtli\nMN8To9HgBzd6EdNTZs+77noZ6lvajJqaRkhs1p6WJNTv/wAWkpDVmDfgAX0drUL6\nc6fw4BynBZiyKxYsnB89SXwGPQsMSHD1wp+I8hWvwJhhNNBhXTHRF2fyl9OY+wbx\nErFnVAOyek06DUG2w7wVFxaVLlYMdaFrpbqCnlCQXJZJVq7+BF0hlv4p2RJMcRzn\nm5AaKgS3gT1VbqqSfJ+qVskS6sMxQFh8mo88NfMYDndp3htVdSH5TOLY6bPZ/40H\nV6IC284JAgMBAAECggEACrUKOD+ttBFr/4kuQshpAYPiXGSWmJueekkFygAAIJsI\nDyoywRlA9AxW51fonoUjWWvL8mfnFan/yY3WJ6Wz5upJQ9F0DQn/RnhD3kyo1CF4\nlOYY2RLx0hnpaz5V5JQNonzro3KgpRMcJIdpaecvHX2bXYixA/1HDTaxMmmF+rPw\n06CSK74+e3yeLBCMQfCG0q7tfyvMdguczD5IVlHOHCGvwKubd6G2DzZvWTDYsUjU\nSB+5+A36CR78WEgEWzo+rZAJbDySzJzaZc7L5PD/M2knZuHQtKFjIRLHhNBCDP1q\nXGPwJsX3UW8FdZ3hvNTQg22SrqujgvwSMDke2H18IQKBgQDHhwSoeMCKFUen4QuY\nJBeNp/nHmb5iSDN53OF9FLuc+u3OZyCz8loQbjzIQRFGyp/5dSK1p5MhH/BjWZ3l\na4fFcuzzkz6zyoDTqbjLv7GpLiVK8wIwfrahO04fEbEmnM+EDKv35Uo6ht08nu4b\nsmgJ1z8DPunVXlhGybIAbSrnUQKBgQC6X68ZB2asRm6frOOwHLE31H9uIVKbJjHu\naOKmZlP2o2wELWDyjIIQZahKMDnyinLAug8Cscp2OHJf4zWYBfR00pHFXSUGkVND\nycC1Lz8GIO3jdvyx6LzMh2OtFjWNMtCLbAjLXkfb8iPJw0ZpqpLRt8U9Uk7tHbRG\nc1A6XDA9OQKBgCD0608MivkD5M8U+/5IT9+lFDvk6C6BsIb7df9cElUumWMTY7J1\nYG0AWGfXX4wq4dupfm8027eH+APhBJSle0qg3gSpmJzH4RmVGiIFasoABkbn9r+d\n3nqpOhElsfYnxpsQIMOUivs51Ycy1S+b+1VMyWq21Jbau4gNvqoVXhXRAoGAJH6S\nCeOiHj/Yb5nqJ9Umepk4rrcFtu2+v0F4iD7nWBdeEl9UaYpL+av+TTCuWCj2GXkV\ncWChFY8uDkqudutLmAiXlL8NfgC8/jwmaRQsUiXmjzEAgFHjjmVAhmcf61s07Ogl\nvLTke1Qp39tGEXDeOQS0MbLJU7MKVvVDk3nz1DkCgYA1DH1FJQa2M+AmdOlRYC09\nomyx+NIdyNbISIY4ZQfrCiig62Yh4p0i1Addk0ZMNNnV3kI9jxKnNOk7RhDQW8M8\nLKrGsrXcIvJvx1FI/HHqVCOsJ7q0971J8lafBDF92CXEg9Jjv2dnsUXXMbSZWmQJ\nLpNkIpv3s9dmICpoWJ+4hQ==\n-----END PRIVATE KEY-----\n",
		"client_email": "wopta-dev-frontend-sa@positive-apex-350507.iam.gserviceaccount.com",
		"client_id": "109633118059069778723",
		"auth_uri": "https://accounts.google.com/o/oauth2/auth",
		"token_uri": "https://oauth2.googleapis.com/token",
		"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
		"client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/wopta-dev-frontend-sa%40positive-apex-350507.iam.gserviceaccount.com"
	  } `)
	enrich_vat := r.PathPrefix("/enrich-vat").Subrouter()
	rules_sub := r.PathPrefix("/rules").Subrouter()
	quote := r.PathPrefix("/quote").Subrouter()

	enrich_vat.HandleFunc("/{key}/{key}", enr.EnrichVat).Methods("GET")
	rules_sub.HandleFunc("/{key}", rules.Rules).Methods("POST")
	quote.HandleFunc("/quote", q.QuoteAllrisk).Methods("POST")
	http.Handle("/", r)
	fmt.Println("Listening on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
