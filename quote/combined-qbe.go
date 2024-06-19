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
		policy     *models.Policy
		inputCells []Cell
	)

	log.SetPrefix("[CombinedQbeFx] ")
	defer log.SetPrefix("")

	log.Println("Handler start -----------------------------------------------")

	req := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	err := json.Unmarshal(req, &policy)
	lib.CheckError(err)
	inputCells = append(inputCells, setInputCell(policy)...)
	qs := QuoteSpreadsheet{
		Id:          "1tn0Jqce-r_JKdecExFOFVEJdGUaPYdGo31A9FOgvt-Y",
		InputCells:  inputCells,
		OutputCells: setOutputCell(),
		InitCells:   setInitCells(),
		SheetName:   "Input dati Polizza",
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
func setInputCell(policy *models.Policy) []Cell {
	var inputCells []Cell
	assEnterprise := getAssetByType(policy, "enterprise")
	assBuildings := getAssetByType(policy, "building")

	inputCells = append(inputCells, Cell{Cell: "C10", Value: policy.StartDate.Format("02-01-2006")})
	inputCells = append(inputCells, Cell{Cell: "C11", Value: policy.EndDate.Format("02-01-2006")})
	inputCells = append(inputCells, setEnterpriseCell(assEnterprise[0])...)
	for _, eg := range assEnterprise[0].Guarantees {
		inputCells = append(inputCells, getEnterpriseGuaranteCellsBySlug(eg)...)
	}
	for i, build := range assBuildings {
		for _, bg := range build.Guarantees {
			inputCells = append(inputCells, getBuildingGuaranteCellsBySlug(bg, i)...)
		}

	}
	return inputCells
}
func setEnterpriseCell(assets models.Asset) []Cell {
	var inputCells []Cell

	inputCells = append(inputCells, Cell{Cell: "E6", Value: assets.Enterprise.VatCode})
	inputCells = append(inputCells, Cell{Cell: "C5", Value: assets.Enterprise.Name})

	return inputCells
}

func getAssetByType(policy *models.Policy, asstype string) []models.Asset {
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
func getAssetGuarante(assets *models.Asset, slug string) models.Guarante {
	var (
		guarante models.Guarante
	)
	for _, g := range assets.Guarantees {
		if g.Slug == slug {
			guarante = g
		}

	}
	return guarante
}
func getEnterpriseGuaranteCellsBySlug(guarante models.Guarante) []Cell {
	var (
		cells []Cell
	)
	switch guarante.Slug {

	case "electrical-phenomenon":
		cells = []Cell{
			{
				Cell:  "C48",
				Value: guarante.SumInsuredLimitOfIndemnity,
			},
		}
	case "refrigeration-goods":
		cells = []Cell{
			{
				Cell:  "C49",
				Value: guarante.SumInsuredLimitOfIndemnity,
			},
		}
	case "machinery-breakdown":
		cells = []Cell{
			{
				Cell:  "C50",
				Value: guarante.SumInsuredLimitOfIndemnity,
			},
		}
	case "electronic-equipment":
		cells = []Cell{
			{
				Cell:  "C51",
				Value: guarante.SumInsuredLimitOfIndemnity,
			},
		}
	case "theft":
		cells = []Cell{
			{
				Cell:  "C52",
				Value: guarante.SumInsuredLimitOfIndemnity,
			},
		}
	case "third-party-recourse":
		cells = []Cell{
			{
				Cell:  "C47",
				Value: guarante.SumInsuredLimitOfIndemnity,
			},
		}
	case "third-party-liability-work-providers":
		cells = []Cell{
			{
				Cell:  "C66",
				Value: guarante.SumInsuredLimitOfIndemnity,
			},
		}
	case "product-liability":
		cells = []Cell{
			{
				Cell:  "C67",
				Value: guarante.SumInsuredLimitOfIndemnity,
			},
		}
	case "management-organization":
		cells = []Cell{
			{
				Cell:  "C68",
				Value: guarante.SumInsuredLimitOfIndemnity,
			},
		}
	case "cyber":
		cells = []Cell{
			{
				Cell:  "C74",
				Value: guarante.SumInsuredLimitOfIndemnity,
			},
		}
	case "daily-allowance":
		cells = []Cell{
			{
				Cell:  "C58",
				Value: guarante.SumInsuredLimitOfIndemnity,
			},
		}
	case "increased-cost":
		cells = []Cell{
			{
				Cell:  "C59",
				Value: guarante.SumInsuredLimitOfIndemnity,
			},
		}
	case "loss-rent":
		cells = []Cell{
			{
				Cell:  "C61",
				Value: guarante.SumInsuredLimitOfIndemnity,
			},
		}
	}
	return cells
}
func getBuildingGuaranteCellsBySlug(guarante models.Guarante, colum int) []Cell {
	var (
		cells []Cell
	)
	col := map[int]string{0: "C", 1: "D", 2: "E", 3: "F", 4: "G"}
	switch guarante.Slug {
	case "building":
		cells = []Cell{
			{
				Cell:  col[colum] + "41",
				Value: guarante.Value.SumInsured,
			},
		}

	case "rental-risk":
		cells = []Cell{
			{
				Cell:  col[colum] + "42",
				Value: guarante.SumInsuredLimitOfIndemnity,
			},
		}
	case "machinery":
		cells = []Cell{
			{
				Cell:  col[colum] + "43",
				Value: guarante.SumInsuredLimitOfIndemnity,
			},
		}
	case "goods":
		cells = []Cell{
			{
				Cell:  col[colum] + "44",
				Value: guarante.SumInsuredLimitOfIndemnity,
			},
		}
	case "goods-temporary-increase":
		cells = []Cell{
			{
				Cell:  col[colum] + "45",
				Value: guarante.SumInsuredLimitOfIndemnity,
			},
		}

	}
	return cells
}
func setInitCells() []Cell {

	cells := []Cell{
		{
			Cell:  "C4",
			Value: "0",
		}, {
			Cell:  "C5",
			Value: "0",
		}, {
			Cell:  "E6",
			Value: "0",
		}, {
			Cell:  "G6",
			Value: "0",
		}, {
			Cell:  "C7",
			Value: "0",
		}, {
			Cell:  "C8",
			Value: "0",
		}, {
			Cell:  "C9",
			Value: "0",
		}, {
			Cell:  "D9",
			Value: "0",
		}, {
			Cell:  "G9",
			Value: "0",
		}, {
			Cell:  "G10",
			Value: "0",
		}, {
			Cell:  "G11",
			Value: "NO",
		}, {
			Cell:  "C10",
			Value: "0",
		}, {
			Cell:  "C11",
			Value: "0",
		}, {
			Cell:  "C12",
			Value: "0",
		}, {
			Cell:  "C19",
			Value: "",
		}, {
			Cell:  "C20",
			Value: "",
		}, {
			Cell:  "D19",
			Value: "0",
		}, {
			Cell:  "D20",
			Value: "",
		}, {
			Cell:  "E19",
			Value: "",
		}, {
			Cell:  "E20",
			Value: "",
		}, {
			Cell:  "F19",
			Value: "",
		}, {
			Cell:  "F20",
			Value: "",
		}, {
			Cell:  "G19",
			Value: "",
		}, {
			Cell:  "G20",
			Value: "",
		}, {
			Cell:  "C21",
			Value: "",
		}, {
			Cell:  "C29",
			Value: "0",
		}, {
			Cell:  "C30",
			Value: "0",
		}, {
			Cell:  "C31",
			Value: "0",
		}, {
			Cell:  "C32",
			Value: "0",
		}, {
			Cell:  "C33",
			Value: "0",
		}, {
			Cell:  "C34",
			Value: "NO",
		}, {
			Cell:  "C35",
			Value: "NO",
		}, {
			Cell:  "C36",
			Value: "NO",
		}, {
			Cell:  "C41",
			Value: "0",
		}, {
			Cell:  "C42",
			Value: "0",
		}, {
			Cell:  "C43",
			Value: "0",
		}, {
			Cell:  "C45",
			Value: "0",
		},
		{
			Cell:  "D21",
			Value: "",
		}, {
			Cell:  "D29",
			Value: "0",
		}, {
			Cell:  "D30",
			Value: "0",
		}, {
			Cell:  "D31",
			Value: "0",
		}, {
			Cell:  "D32",
			Value: "0",
		}, {
			Cell:  "D33",
			Value: "0",
		}, {
			Cell:  "D34",
			Value: "NO",
		}, {
			Cell:  "D35",
			Value: "NO",
		}, {
			Cell:  "D36",
			Value: "NO",
		}, {
			Cell:  "D41",
			Value: "0",
		}, {
			Cell:  "D42",
			Value: "0",
		}, {
			Cell:  "D43",
			Value: "0",
		}, {
			Cell:  "D45",
			Value: "0",
		},

		{
			Cell:  "E21",
			Value: "",
		}, {
			Cell:  "E29",
			Value: "0",
		}, {
			Cell:  "E30",
			Value: "0",
		}, {
			Cell:  "E31",
			Value: "0",
		}, {
			Cell:  "E32",
			Value: "0",
		}, {
			Cell:  "E33",
			Value: "0",
		}, {
			Cell:  "E34",
			Value: "NO",
		}, {
			Cell:  "E35",
			Value: "NO",
		}, {
			Cell:  "E36",
			Value: "NO",
		}, {
			Cell:  "E41",
			Value: "0",
		}, {
			Cell:  "E42",
			Value: "0",
		}, {
			Cell:  "E43",
			Value: "0",
		}, {
			Cell:  "E45",
			Value: "0",
		},

		{
			Cell:  "F21",
			Value: "0",
		}, {
			Cell:  "F29",
			Value: "0",
		}, {
			Cell:  "F30",
			Value: "0",
		}, {
			Cell:  "F31",
			Value: "0",
		}, {
			Cell:  "F32",
			Value: "0",
		}, {
			Cell:  "F33",
			Value: "0",
		}, {
			Cell:  "F34",
			Value: "NO",
		}, {
			Cell:  "F35",
			Value: "NO",
		}, {
			Cell:  "F36",
			Value: "NO",
		}, {
			Cell:  "F41",
			Value: "0",
		}, {
			Cell:  "F42",
			Value: "0",
		}, {
			Cell:  "F43",
			Value: "0",
		}, {
			Cell:  "F45",
			Value: "0",
		},

		{
			Cell:  "G21",
			Value: "0",
		}, {
			Cell:  "G29",
			Value: "0",
		}, {
			Cell:  "G30",
			Value: "0",
		}, {
			Cell:  "G31",
			Value: "0",
		}, {
			Cell:  "G32",
			Value: "0",
		}, {
			Cell:  "G33",
			Value: "0",
		}, {
			Cell:  "G34",
			Value: "NO",
		}, {
			Cell:  "G35",
			Value: "NO",
		}, {
			Cell:  "G36",
			Value: "NO",
		}, {
			Cell:  "G41",
			Value: "0",
		}, {
			Cell:  "G42",
			Value: "0",
		}, {
			Cell:  "G43",
			Value: "0",
		}, {
			Cell:  "G45",
			Value: "0",
		},
		{
			Cell:  "C46",
			Value: "0",
		}, {
			Cell:  "C47",
			Value: "0",
		}, {
			Cell:  "C48",
			Value: "0",
		}, {
			Cell:  "C49",
			Value: "0",
		}, {
			Cell:  "C50",
			Value: "0",
		}, {
			Cell:  "C51",
			Value: "0",
		}, {
			Cell:  "C52",
			Value: "0",
		}, {
			Cell:  "C57",
			Value: "Esclusa",
		}, {
			Cell:  "C61",
			Value: "0",
		}, {
			Cell:  "C67",
			Value: "NO",
		}, {
			Cell:  "C68",
			Value: "NO",
		}, {
			Cell:  "C69",
			Value: "0",
		}, {
			Cell:  "C70",
			Value: "0",
		}, {
			Cell:  "F66",
			Value: "0",
		}, {
			Cell:  "F67",
			Value: "0",
		}, {
			Cell:  "F68",
			Value: "0",
		},
		{
			Cell:  "F69",
			Value: "NO",
		}, {
			Cell:  "F70",
			Value: "0",
		},
	}
	return cells
}
