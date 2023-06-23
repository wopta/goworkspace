package models

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

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
	Payload  string    `bigquery:"payload"`
	Date     time.Time `bigquery:"date"`
	UserUid  string    `bigquery:"userUid"`
	Method   string    `bigquery:"method"`
	Endpoint string    `bigquery:"endpoint"`
	Role     string    `bigquery:"role"`
}

func ParseHttpRequest(r *http.Request) (AuditLog, error) {
	if r.Body != nil {
		return AuditLog{}, fmt.Errorf("no payload found: body is empty!")
	}

	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return AuditLog{}, fmt.Errorf("cannot retrieve the payload: %v", err)
	}

	userUid := r.Header.Get("userUid")
	if userUid == "" {
		return AuditLog{}, fmt.Errorf("cannot retrieve the user's UID. Header (userUid) is empty")
	}

	role, err := getRole(r)
	if err != nil {
		return AuditLog{}, fmt.Errorf("cannot retrieve the role: %s", err)
	}

	return AuditLog{
		Payload:  string(body),
		Date:     time.Now().UTC(),
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

// Copy and paseted from router/router.go: VerifyAuthorization.
// Modified so that it returns an error in case.
func getRole(r *http.Request) (string, error) {
	idToken := strings.ReplaceAll(r.Header.Get("Authorization"), "Bearer ", "")
	if idToken == "" {
		return "", fmt.Errorf("the authorization Header is empty!")
	}

	token, err := lib.VerifyUserIdToken(idToken)
	if err != nil {
		return "", fmt.Errorf("verify ID token error: %s", err)
	}
	if token.Claims["role"] == nil {
		return "", fmt.Errorf("user role not set")
	}
	role := token.Claims["role"].(string)
	return role, nil
}
