package models

import (
	"encoding/json"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/wopta/goworkspace/lib"
	"google.golang.org/api/iterator"
)

func UnmarshalProduct(data []byte) (Product, error) {
	var r Product
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Product) Marshal() ([]byte, error) {

	return json.Marshal(r)
}

type Product struct {
	NameDesc       *string          `firestore:"nameDesc,omitempty" json:"nameDesc,omitempty"`
	Companies      []Company        `firestore:"companies,omitempty" json:"companies,omitempty"`
	ProductUid     string           `firestore:"productUid,omitempty" json:"productUid,omitempty"`
	ProductVersion int              `firestore:"productVersion,omitempty" json:"productVersion,omitempty"`
	Version        string           `firestore:"version,omitempty" json:"version,omitempty"`
	Number         int              `firestore:"number,omitempty" json:"number,omitempty"`
	Name           string           `firestore:"name,omitempty" json:"name,omitempty"`
	Steps          []Step           `firestore:"steps,omitempty" json:"steps"`
	Offers         map[string]Offer `firestore:"offers,omitempty" json:"offers,omitempty"`
}

type Company struct {
	Name                string              `firestore:"name,omitempty" json:"name,omitempty"`
	Code                string              `firestore:"code,omitempty" json:"code,omitempty"`
	Commission          float64             `firestore:"commission,omitempty" json:"commission,omitempty"`
	CommissionRenew     float64             `firestore:"commissionRenew,omitempty" json:"commissionRenew,omitempty"`
	MinimumMonthlyPrice float64             `firestore:"minimumMonthlyPrice,omitempty" json:"minimumMonthlyPrice,omitempty"`
	MinimumYearlyPrice  float64             `firestore:"minimumYearlyPrice,omitempty" json:"minimumYearlyPrice,omitempty"`
	Guarantees          *[]Guarante         `firestore:"guarantees,omitempty" json:"guarantees,omitempty"`
	GuaranteesMap       map[string]Guarante `firestore:"guaranteesMap,omitempty" json:"guaranteesMap,omitempty"`
}
type Step struct {
	Widget     string      `firestore:"widget,omitempty" json:"widget"`
	Attributes interface{} `firestore:"attributes,omitempty" json:"attributes"`
	Children   []Child     `firestore:"children,omitempty" json:"children,omitempty"`
}

type Child struct {
	Widget     string      `firestore:"widget,omitempty" json:"widget"`
	Attributes interface{} `firestore:"attributes,omitempty" json:"attributes"`
}

type Offer struct {
	Name        string `firestore:"name,omitempty" json:"name,omitempty"`
	Description string `firestore:"description,omitempty" json:"description,omitempty"`
	Order       int    `firestore:"order,omitempty" json:"order,omitempty"`
}

func ProductToListData(query *firestore.DocumentIterator) []Product {
	var result []Product
	for {
		d, err := query.Next()
		if err != nil {
			if err == iterator.Done {

				break
			}
		}
		var value Product
		e := d.DataTo(&value)
		lib.CheckError(e)
		result = append(result, value)

		log.Println(len(result))
	}
	return result
}
