package models

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"cloud.google.com/go/civil"
	"github.com/wopta/goworkspace/lib"
)

const (
	auditTable = "audit"
	dataset    = "wopta"
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
	token, err := lib.VerifyUserIdToken(idToken)
	if err != nil {
		return AuditLog{}, fmt.Errorf("cannot retrieve the user's token: %v", err)
	}

	userUid, err := lib.GetUserIdFromIdToken(idToken)
	if err != nil {
		return AuditLog{}, fmt.Errorf("cannot retrieve the user's UID. Header (userUid) is empty")
	}

	role := ""
	if token.Claims["role"] == nil {
		return AuditLog{}, fmt.Errorf("cannot retrieve the role: %s", err)
	}
	role = token.Claims["role"].(string)

	return AuditLog{
		Payload:  body,
		Date:     civil.DateTimeOf(time.Now().UTC()),
		UserUid:  userUid,
		Method:   r.Method,
		Endpoint: r.RequestURI,
		Role:     role,
	}, nil
}

func (a AuditLog) SaveToBigQuery() error {
	if err := lib.InsertRowsBigQuery(dataset, auditTable, a); err != nil {
		return fmt.Errorf("cannot save the audit-log: %v", err)
	}
	return nil
}
