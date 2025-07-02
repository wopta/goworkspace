package models

import (
	"strconv"

	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
)

// {"responses":{"TIPO MOVIMENTO":"Inserimento","Targa Inserimento":"test","MODELLO VEICOLO":"test mod","DATA IMMATRICOLAZIONE":"1212-12-02","DATA INIZIO VALIDITA' COPERTURA":"1212-12-12"}}
// {"responses":{"TIPO MOVIMENTO":"Annullo","Targa Annullo":"targa","DATA FINE VALIDITA' COPERTURA":"0009-09-09"},"mail":"test@gmail.com"}
type FirestoreSequence struct {
	Last   int    `json:"last" firestore:"last" bigquery:"last"`
	Name   string `json:"name" firestore:"name" bigquery:"name"`
	Prefix string `json:"prefix" firestore:"prefix" bigquery:"prefix"`
	Format string `json:"format" firestore:"format" bigquery:"format"` //`%07d`
	Start  string `json:"start" firestore:"start" bigquery:"start"`
}

func GetFirestoreSequenceLast(name string) string {
	q := lib.FireGenericQueries[FirestoreSequence]{
		Queries: []lib.Firequery{
			{
				Field:      "name",
				Operator:   "==",
				QueryValue: name,
			},
		}}

	node, uid, e := q.FireQueryUid("sequence")
	if e != nil {
		log.Error(e)
	}
	if len(node) > 0 {
		last := node[0].Add()
		node[0].Last = last
		lib.SetFirestoreErr("sequence", uid[0], node[0])
		return strconv.Itoa(last)
	}
	return ""
}
func (s *FirestoreSequence) Add() int {
	result := s.Last + 1
	s.Last = result
	return result
}
