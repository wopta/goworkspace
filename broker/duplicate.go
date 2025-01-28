package broker

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/mohae/deepcopy"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	plc "github.com/wopta/goworkspace/policy"
)

var (
	errOperationNotAllowed = errors.New("operation not allowed")
)

const (
	statusDuplicated = "Duplicated"
)

func DuplicateFx(_ http.ResponseWriter, r *http.Request) (string, any, error) {
	var (
		err            error
		originalPolicy models.Policy
		responseBytes  []byte
	)

	log.SetPrefix("[AcceptanceFx]")
	defer func() {
		if err != nil {
			log.Printf("error: %s", err.Error())
		}
		log.Println("Handler end ---------------------------------------------")
		log.SetPrefix("")
	}()
	log.Println("Handler start -----------------------------------------------")

	policyUid := chi.URLParam(r, "uid")

	if originalPolicy, err = plc.GetPolicy(policyUid, ""); err != nil {
		log.Println("could not retrieve policy from DB")
		return "", nil, err
	}

	if originalPolicy.CompanyEmit {
		log.Println("cannot duplicate already emitted policy")
		err = errOperationNotAllowed
		return "", nil, err
	}

	duplicatedPolicy := deepcopy.Copy(originalPolicy).(models.Policy)

	now := time.Now().UTC()

	originalPolicy.Status = models.PolicyStatusDeleted
	originalPolicy.StatusHistory = append(originalPolicy.StatusHistory, statusDuplicated, originalPolicy.Status)
	originalPolicy.DeleteDesc = "Annulata per modifica"
	originalPolicy.DeleteDate = now
	originalPolicy.Updated = now

	duplicatedPolicy.Uid = lib.NewDoc(lib.PolicyCollection)
	duplicatedPolicy.ProposalNumber = 0
	duplicatedPolicy.Attachments = nil
	duplicatedPolicy.Status = models.PolicyStatusInit
	duplicatedPolicy.StatusHistory = []string{statusDuplicated, duplicatedPolicy.Status}
	duplicatedPolicy.Updated = now

	if err = lib.SetFirestoreErr(lib.PolicyCollection, originalPolicy.Uid, originalPolicy); err != nil {
		log.Println("error updating original policy")
		return "", nil, err
	}
	originalPolicy.BigquerySave("")

	if err = lib.SetFirestoreErr(lib.PolicyCollection, duplicatedPolicy.Uid, duplicatedPolicy); err != nil {
		log.Println("error updating duplicated policy")
		return "", nil, err
	}
	duplicatedPolicy.BigquerySave("")

	if responseBytes, err = duplicatedPolicy.Marshal(); err != nil {
		log.Println("error marshalling response")
		return "", nil, err
	}

	return string(responseBytes), duplicatedPolicy, err
}
