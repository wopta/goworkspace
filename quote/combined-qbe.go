package quote

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func CombinedQbeFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		policy *models.Policy
	)

	log.SetPrefix("[CombinedQbeFx] ")
	defer log.SetPrefix("")

	log.Println("Handler start -----------------------------------------------")

	req := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	err := json.Unmarshal(req, &policy)
	lib.CheckError(err)
	inputCells := []Cell{{
		Cell: "",
	},
	}
	qs := QuoteSpreadsheet{Id: "1GMtY4EIR2qeyylTOoCfNLFWVNam0H6MF1Is8yD2DiWI",
		InputCells:  inputCells,
		OutputCells: setOutputCell(),
	}
	qs.Spreadsheets()

	policyJson, err := policy.Marshal()

	log.Println("Handler end -------------------------------------------------")

	return string(policyJson), policy, err
}
func setOutputCell() []Cell {

	res := []Cell{{
		Cell: "C81",
	}, {
		Cell: "C82",
	}, {
		Cell: "C83",
	}, {
		Cell: "C84",
	}, {
		Cell: "C85",
	}, {
		Cell: "C86",
	}, {
		Cell: "C87",
	}, {
		Cell: "C88",
	}, {
		Cell: "C89",
	}, {
		Cell: "C90",
	}, {
		Cell: "C91",
	}, {
		Cell: "C92",
	}, {
		Cell: "C93",
	}, {
		Cell: "C94",
	}, {
		Cell: "C95",
	}, {
		Cell: "C96",
	}, {
		Cell: "C97",
	}, {
		Cell: "C98",
	}, {
		Cell: "C99",
	}, {
		Cell: "C100",
	}, {
		Cell: "C101",
	},
	}

	return res
}
func setInputCell(policy models.Policy) []Cell {
	var res []Cell

	return res
}
func getAssetByType(policy models.Policy, asstype string) []models.Asset {
	var (
		assets []models.Asset
	)
	for _, asset := range policy.Assets {
		if asset.Type == asstype {
			assets = append(assets, asset)
		}

	}
	return assets
}
func getAssetGuarante(assets models.Asset, slug string) models.Guarante {
	var (
		guarante  models.Guarante 
	)
	for _, g := range assets.Guarantees {
		if g.Slug == slug {
			guarante = g
		}

	}
	return guarante
}
