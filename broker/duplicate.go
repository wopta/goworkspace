package broker

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
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
		policy models.Policy
		responseBytes  []byte
	)

	log.SetPrefix("[DuplicateFx]")
	defer func() {
		if err != nil {
			log.Printf("error: %s", err.Error())
		}
		log.Println("Handler end ---------------------------------------------")
		log.SetPrefix("")
	}()
	log.Println("Handler start -----------------------------------------------")

	policyUid := chi.URLParam(r, "uid")

	if policy, err = plc.GetPolicy(policyUid, ""); err != nil {
		log.Println("could not retrieve policy from DB")
		return "", nil, err
	}

	if policy.CompanyEmit {
		log.Println("cannot duplicate already emitted policy")
		err = errOperationNotAllowed
		return "", nil, err
	}

	now := time.Now().UTC()

	policy.Uid = lib.NewDoc(lib.PolicyCollection)
	policy.ProposalNumber = 0
	policy.Attachments = nil
	policy.Status = models.PolicyStatusInit
	policy.StatusHistory = []string{statusDuplicated, policy.Status}
	policy.CreationDate = now
	policy.Updated = now

	if err = lib.SetFirestoreErr(lib.PolicyCollection, policy.Uid, policy); err != nil {
		log.Println("error updating duplicated policy")
		return "", nil, err
	}
	policy.BigquerySave("")

	if responseBytes, err = policy.Marshal(); err != nil {
		log.Println("error marshalling response")
		return "", nil, err
	}

	return string(responseBytes), policy, err
}
