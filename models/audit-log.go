package models

import (
	"fmt"
	"io"
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

func ParseHttpRequest(r *http.Request) (AuditLog, error) {
	body := ""
	if r.Body != nil {
		defer r.Body.Close()
		body_bytes, err := io.ReadAll(r.Body)
		if err != nil {
			return AuditLog{}, fmt.Errorf("cannot retrieve the payload: %v", err)
		}
		body = string(body_bytes)
	}

	idToken := r.Header.Get("Authorization")
	authToken, err := GetAuthTokenFromIdToken(idToken)
	if err != nil {
		return AuditLog{}, fmt.Errorf("cannot retrieve the user's authorization token: %v", err)
	}

	return AuditLog{
		Payload:  body,
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
