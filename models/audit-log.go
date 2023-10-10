package models

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/civil"
	"github.com/wopta/goworkspace/lib"
)

// The AuditLog is a data structure that stores request's data.
// In particular, the Payload should contain data from the request's body.
//
// Example:
//
//	a := ParseHttpRequest(r)
//	a.SaveToBigQuery()
type AuditLog struct {
	Payload  string         `bigquery:"payload"`
	Date     civil.DateTime `bigquery:"date"`
	UserUid  string         `bigquery:"userUid"`
	Method   string         `bigquery:"method"`
	Endpoint string         `bigquery:"endpoint"`
	Role     string         `bigquery:"role"`
}

// Since io.ReadAll is a stream we cannot read r.Body twice
// so we pass the string that was read by the original
// function as a param and use it in the AuditLog.Payload
func ParseHttpRequest(r *http.Request, payload string) (AuditLog, error) {
	idToken := r.Header.Get("Authorization")
	authToken, err := GetAuthTokenFromIdToken(idToken)
	if err != nil {
		return AuditLog{}, fmt.Errorf("cannot retrieve the user's authorization token: %v", err)
	}

	return AuditLog{
		Payload:  payload,
		Date:     civil.DateTimeOf(time.Now().UTC()),
		UserUid:  authToken.UserID,
		Method:   r.Method,
		Endpoint: r.RequestURI,
		Role:     authToken.Role,
	}, nil
}

func (a AuditLog) SaveToBigQuery() error {
	if err := lib.InsertRowsBigQuery(WoptaDataset, AuditsCollection, a); err != nil {
		return fmt.Errorf("cannot save the audit-log: %v", err)
	}
	return nil
}

func CreateAuditLog(r *http.Request, payload string) {
	log.Println("[CreateAuditLog] saving audit trail...")
	audit, err := ParseHttpRequest(r, payload)
	if err != nil {
		log.Printf("[CreateAuditLog] error creating audit log: %s", err.Error())
	}
	log.Printf("[CreateAuditLog] audit log: %v", audit)
	audit.SaveToBigQuery()
}
