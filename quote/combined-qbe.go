package quote

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

const (
	dateFormat = "02/01/2006"
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
	log.Println("Request: ", string(req))
	err := json.Unmarshal(req, &policy)
	lib.CheckError(err)
	b, err := json.Marshal(policy)
	log.Println("Request Marshal: ", string(b))
	lib.CheckError(err)
	inputCells = append(inputCells, setInputCell(policy)...)
	qs := QuoteSpreadsheet{
		Id:          "1tn0Jqce-r_JKdecExFOFVEJdGUaPYdGo31A9FOgvt-Y",
		InputCells:  inputCells,
		OutputCells: setOutputCell(),
		InitCells:   setInitCells(),
		SheetName:   "Input dati Polizza",
	}
	outCells := qs.Spreadsheets()
	mapCellPolicy(policy, outCells)

	policyJson, err := policy.Marshal()
	log.Println("Response: ", string(policyJson))
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
	},
	}

	return res
}
func mapCellPolicy(policy *models.Policy, cells []Cell) {
	var priceGroup []models.Price

	policy.IsReserved = true
	policy.Channel = "network"

	policy.OffersPrices = map[string]map[string]*models.Price{
		"default": {
			"yearly":  &models.Price{},
			"monthly": &models.Price{},
		},
	}
	for _, cell := range cells {
		s, err := strconv.ParseFloat(strings.Trim(strings.Replace(strings.Replace(cell.Value.(string), ".", "", -1), ",", ".", -1), " "), 64)
		log.Println(err)
		switch cell.Cell {
		case "C81":
			if err == nil {
				priceGroup = append(priceGroup, models.Price{
					Name: "Fabbricato",
					Net:  s,
				})
			} else {
				priceGroup = append(priceGroup, models.Price{
					Name:        "Fabbricato",
					Description: cell.Value.(string),
				})
			}

		case "C82":
			if err == nil {
				priceGroup = append(priceGroup, models.Price{
					Name: "Contenuto (Merci e Macchinari)",
					Net:  s,
				})
			} else {
				priceGroup = append(priceGroup, models.Price{
					Name:        "Contenuto (Merci e Macchinari)",
					Description: cell.Value.(string),
				})
			}

		case "C83":
			if err == nil {
				priceGroup = append(priceGroup, models.Price{
					Name: "Merci (aumento temporaneo)",
					Net:  s,
				})
			} else {
				priceGroup = append(priceGroup, models.Price{
					Name:        "Merci (aumento temporaneo)",
					Description: cell.Value.(string),
				})
			}

		case "C84":
			if err == nil {
				priceGroup = append(priceGroup, models.Price{
					Name: "Furto, rapina, estorsione (in aumento)",
					Net:  s,
				})
			} else {
				priceGroup = append(priceGroup, models.Price{
					Name:        "Furto, rapina, estorsione (in aumento)",
					Description: cell.Value.(string),
				})
			}

		case "C85":
			if err == nil {
				priceGroup = append(priceGroup, models.Price{
					Name: "Rischio locativo (in aumento)",
					Net:  s,
				})
			} else {
				priceGroup = append(priceGroup, models.Price{
					Name:        "Rischio locativo (in aumento)",
					Description: cell.Value.(string),
				})
			}

		case "C86":

			log.Println(err)
			if err == nil {
				priceGroup = append(priceGroup, models.Price{
					Name: "Altre garanzie su Contenuto",
					Net:  s,
				})
			} else {
				priceGroup = append(priceGroup, models.Price{
					Name:        "Altre garanzie su Contenuto",
					Description: cell.Value.(string),
				})
			}

		case "C87":
			if err == nil {
				priceGroup = append(priceGroup, models.Price{
					Name: "Ricorso terzi (in aumento)",
					Net:  s,
				})
			} else {
				priceGroup = append(priceGroup, models.Price{
					Name:        "Ricorso terzi (in aumento)",
					Description: cell.Value.(string),
				})
			}

		case "C88":
			if err == nil {
				priceGroup = append(priceGroup, models.Price{
					Name: "Danni indiretti",
					Net:  s,
				})
			} else {
				priceGroup = append(priceGroup, models.Price{
					Name:        "Danni indiretti",
					Description: cell.Value.(string),
				})
			}

		case "C89":

			log.Println(err)
			priceGroup = append(priceGroup, models.Price{
				Name: "Perdita Pigioni",
				Net:  s,
			})
		case "C90":
			if err == nil {
				priceGroup = append(priceGroup, models.Price{
					Name: "Responsabilità civile terzi",
					Net:  s,
				})
			} else {
				priceGroup = append(priceGroup, models.Price{
					Name:        "Responsabilità civile terzi",
					Description: cell.Value.(string),
				})
			}

		case "C91":
			if err == nil {
				priceGroup = append(priceGroup, models.Price{
					Name: "Responsabilità civile prestatori lavoro",
					Net:  s,
				})
			} else {
				priceGroup = append(priceGroup, models.Price{
					Name:        "Responsabilità civile prestatori lavoro",
					Description: cell.Value.(string),
				})
			}

		case "C92":
			if err == nil {
				priceGroup = append(priceGroup, models.Price{
					Name: "Responsabilità civile prodotti",
					Net:  s,
				})
			} else {
				priceGroup = append(priceGroup, models.Price{
					Name:        "Responsabilità civile prodotti",
					Description: cell.Value.(string),
				})
			}

		case "C93":
			if err == nil {
				priceGroup = append(priceGroup, models.Price{
					Name: "Ritiro Prodotti",
					Net:  s,
				})
			} else {
				priceGroup = append(priceGroup, models.Price{
					Name:        "Ritiro Prodotti",
					Description: cell.Value.(string),
				})
			}

		case "C94":
			if err == nil {
				priceGroup = append(priceGroup, models.Price{
					Name: "Resp. Amministratori Sindaci Dirigenti (D&O)",
					Net:  s,
				})
			} else {
				priceGroup = append(priceGroup, models.Price{
					Name:        "Resp. Amministratori Sindaci Dirigenti (D&O)",
					Description: cell.Value.(string),
				})
			}

		case "C95":
			if err == nil {
				priceGroup = append(priceGroup, models.Price{
					Name: "Cyber",
					Net:  s,
				})
			} else {
				priceGroup = append(priceGroup, models.Price{
					Name:        "Cyber",
					Description: cell.Value.(string),
				})
			}

		case "C96":
			if err == nil {

				policy.OffersPrices["default"]["yearly"].Net = s
				policy.PriceNett = s
			}

		case "C97":
			if err == nil {
				policy.TaxAmount = s
			}

		case "C98":
			if err == nil {
				policy.OffersPrices["default"]["yearly"].Gross = s
			}

		case "C99":

		case "C100":

		default:

		}
	}
	policy.PriceGroup = priceGroup
}
func setInputCell(policy *models.Policy) []Cell {
	var inputCells []Cell
	paymentSplits := map[string]string{

		"yearly":    "Annuale",
		"semestral": "Semestrale",
	}
	assEnterprise := getAssetByType(policy, "enterprise")
	assBuildings := getAssetByType(policy, "building")

	inputCells = append(inputCells, Cell{Cell: "C10", Value: policy.StartDate.Format(dateFormat)})
	inputCells = append(inputCells, Cell{Cell: "C11", Value: policy.StartDate.AddDate(1, 0, 0).Format(dateFormat)})
	inputCells = append(inputCells, setEnterpriseCell(assEnterprise[0])...)
	inputCells = append(inputCells, Cell{Cell: "C24", Value: assEnterprise[0].Enterprise.Employer})
	inputCells = append(inputCells, Cell{Cell: "C25", Value: assEnterprise[0].Enterprise.WorkEmployersRemuneration})
	inputCells = append(inputCells, Cell{Cell: "C26", Value: assEnterprise[0].Enterprise.TotalBilled})
	inputCells = append(inputCells, Cell{Cell: "C16", Value: paymentSplits[policy.PaymentSplit]})
	if policy.PaymentSplit == "semestral" {
		inputCells = append(inputCells, Cell{Cell: "C16", Value: "Semestrale"})
	}
	if policy.PaymentSplit == "yearly" {
		inputCells = append(inputCells, Cell{Cell: "C16", Value: "Annuale"})
	}
	for _, eg := range assEnterprise[0].Guarantees {

		inputCells = append(inputCells, getEnterpriseGuaranteCellsBySlug(eg)...)
	}
	for i, build := range assBuildings {
		col := map[int]string{0: "C", 1: "D", 2: "E", 3: "F", 4: "G"}
		for _, bg := range build.Guarantees {
			inputCells = append(inputCells, Cell{Cell: col[i] + "29", Value: build.Building.Address.PostalCode})
			inputCells = append(inputCells, Cell{Cell: col[i] + "30", Value: build.Building.Address.Locality})
			inputCells = append(inputCells, Cell{Cell: col[i] + "31", Value: build.Building.Address.City})
			inputCells = append(inputCells, Cell{Cell: col[i] + "32", Value: build.Building.Address.StreetName})
			inputCells = append(inputCells, Cell{Cell: col[i] + "19", Value: build.Building.NaicsCategory})
			inputCells = append(inputCells, Cell{Cell: col[i] + "20", Value: build.Building.NaicsDetail})
			inputCells = append(inputCells, Cell{Cell: col[i] + "21", Value: build.Building.Naics})
			inputCells = append(inputCells, Cell{Cell: col[i] + "33", Value: build.Building.BuildingMaterial})
			if build.Building.IsAllarm {
				inputCells = append(inputCells, Cell{Cell: col[i] + "35", Value: "SI"})
			}
			if build.Building.HasSandwitchPanel {
				inputCells = append(inputCells, Cell{Cell: col[i] + "34", Value: "SI"})
			}
			if build.Building.HasSprinkler {
				inputCells = append(inputCells, Cell{Cell: col[i] + "35", Value: "SI"})
			}
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
				Value: guarante.Value.SumInsuredLimitOfIndemnity,
			},
		}
	case "refrigeration-stock":
		cells = []Cell{
			{
				Cell:  "C49",
				Value: guarante.Value.SumInsuredLimitOfIndemnity,
			},
		}
	case "machinery-breakdown":
		cells = []Cell{
			{
				Cell:  "C50",
				Value: guarante.Value.SumInsuredLimitOfIndemnity,
			},
		}
	case "electronic-equipment":
		cells = []Cell{
			{
				Cell:  "C51",
				Value: guarante.Value.SumInsuredLimitOfIndemnity,
			},
		}
	case "theft":
		cells = []Cell{
			{
				Cell:  "C52",
				Value: guarante.Value.SumInsuredLimitOfIndemnity,
			},
			{
				Cell:  "G81",
				Value: int(guarante.Value.Discount),
			},
		}
	case "third-party-recourse":
		cells = []Cell{
			{
				Cell:  "C47",
				Value: guarante.Value.SumInsuredLimitOfIndemnity,
			},
		}
	case "third-party-liability-work-providers":
		cells = []Cell{
			{
				Cell:  "C66",
				Value: guarante.Value.SumInsuredLimitOfIndemnity,
			},
			{
				Cell:  "F66",
				Value: guarante.Value.RetroactiveDate.Format(dateFormat),
			},
			{
				Cell:  "G82",
				Value: int(guarante.Value.Discount),
			},
		}
	case "product-liability":
		cells = []Cell{
			{
				Cell:  "C67",
				Value: guarante.Value.SumInsuredLimitOfIndemnity,
			},
			{
				Cell:  "F67",
				Value: guarante.Value.RetroactiveDate.Format(dateFormat),
			},
			{
				Cell:  "F68",
				Value: guarante.Value.RetroactiveUsaCanDate.Format(dateFormat),
			},
		}
	case "management-organization":
		cells = []Cell{
			{
				Cell:  "C69",
				Value: guarante.Value.SumInsuredLimitOfIndemnity,
			},
			{
				Cell:  "C70",
				Value: guarante.Value.SumInsured,
			}, {
				Cell:  "C69",
				Value: guarante.Value.LimitOfIndemnity,
			},
			{
				Cell:  "C68",
				Value: "SI",
			},
			{
				Cell:  "F70",
				Value: guarante.Value.StartDate.Format(dateFormat),
			},
		}
	case "cyber":
		cells = []Cell{
			{
				Cell:  "C74",
				Value: guarante.Value.SumInsuredLimitOfIndemnity,
			},
		}
	case "daily-allowance":
		cells = []Cell{
			{
				Cell:  "C58",
				Value: guarante.Value.SumInsuredLimitOfIndemnity,
			}, {
				Cell:  "C57",
				Value: "Diaria Giornaliera",
			},
			{
				Cell:  "E58",
				Value: guarante.Value.Duration.Day,
			},
		}
	case "additional-compensation":
		cells = []Cell{

			{
				Cell:  "C57",
				Value: "Indennità Aggiuntiva (+10%)",
			},
		}
	case "increased-cost":
		cells = []Cell{
			{
				Cell:  "C59",
				Value: guarante.Value.SumInsuredLimitOfIndemnity,
			},
			{
				Cell:  "C57",
				Value: "Maggiori Costi",
			},
		}
	case "loss-rent":
		cells = []Cell{
			{
				Cell:  "C61",
				Value: guarante.Value.SumInsuredLimitOfIndemnity,
			},
		}
	case "management-organization-continuity":
		cells = []Cell{
			{
				Cell:  "F70",
				Value: guarante.Value.StartDate.Format(dateFormat),
			},
		}

	case "product-withdrawal":
		cells = []Cell{
			{
				Cell:  "F69",
				Value: "SI",
			},
			{
				Cell:  "C93",
				Value: guarante.Value.SumInsuredLimitOfIndemnity,
			},
			{
				Cell:  "G69",
				Value: guarante.Value.StartDate.Format(dateFormat),
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
				Value: guarante.Value.SumInsuredLimitOfIndemnity,
			},
			{
				Cell:  "G80",
				Value: int(guarante.Value.Discount),
			},
		}

	case "rental-risk":
		cells = []Cell{
			{
				Cell:  col[colum] + "42",
				Value: guarante.Value.SumInsuredLimitOfIndemnity,
			},
		}
	case "machinery":
		cells = []Cell{
			{
				Cell:  col[colum] + "43",
				Value: guarante.Value.SumInsuredLimitOfIndemnity,
			},
		}
	case "stock":
		cells = []Cell{
			{
				Cell:  col[colum] + "44",
				Value: guarante.Value.SumInsuredLimitOfIndemnity,
			},
		}
	case "stock-temporary-increase":
		cells = []Cell{
			{
				Cell:  col[colum] + "45",
				Value: guarante.Value.SumInsuredLimitOfIndemnity,
			},
		}
		if guarante.Value.StartDateString != "" {
			cells = append(cells, []Cell{

				{
					Cell:  "E46",
					Value: guarante.Value.StartDateString,
				},
				{
					Cell:  "C46",
					Value: guarante.Value.Duration.Day,
				}}...)
		}

	}
	return cells
}
func setInitCells() []Cell {

	cells := []Cell{
		{
			Cell:  "C4",
			Value: "",
		}, {
			Cell:  "C5",
			Value: "",
		}, {
			Cell:  "E6",
			Value: "",
		}, {
			Cell:  "G6",
			Value: "",
		}, {
			Cell:  "C7",
			Value: "",
		}, {
			Cell:  "C8",
			Value: "",
		}, {
			Cell:  "C9",
			Value: "",
		}, {
			Cell:  "D9",
			Value: "",
		}, {
			Cell:  "G9",
			Value: "",
		}, {
			Cell:  "G10",
			Value: "",
		}, {
			Cell:  "G11",
			Value: "NO",
		}, {
			Cell:  "C10",
			Value: "",
		}, {
			Cell:  "C11",
			Value: "",
		}, {
			Cell:  "C12",
			Value: "",
		}, {
			Cell:  "C19",
			Value: "",
		}, {
			Cell:  "C20",
			Value: "",
		},
		{
			Cell:  "D19",
			Value: "",
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
		},
		{
			Cell:  "C24",
			Value: "0",
		}, {
			Cell:  "C25",
			Value: "0",
		},
		{
			Cell:  "C26",
			Value: "0",
		}, {
			Cell:  "C29",
			Value: "",
		}, {
			Cell:  "C30",
			Value: "",
		}, {
			Cell:  "C31",
			Value: "",
		}, {
			Cell:  "C32",
			Value: "",
		}, {
			Cell:  "C33",
			Value: "Sconosciuto",
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
			Value: "",
		}, {
			Cell:  "D30",
			Value: "",
		}, {
			Cell:  "D31",
			Value: "",
		}, {
			Cell:  "D32",
			Value: "",
		}, {
			Cell:  "D33",
			Value: "",
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
			Value: "",
		}, {
			Cell:  "E30",
			Value: "",
		}, {
			Cell:  "E31",
			Value: "",
		}, {
			Cell:  "E32",
			Value: "",
		}, {
			Cell:  "E33",
			Value: "",
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
			Value: "",
		}, {
			Cell:  "F29",
			Value: "",
		}, {
			Cell:  "F30",
			Value: "",
		}, {
			Cell:  "F31",
			Value: "",
		}, {
			Cell:  "F32",
			Value: "",
		}, {
			Cell:  "F33",
			Value: "",
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
			Value: "",
		}, {
			Cell:  "G29",
			Value: "",
		}, {
			Cell:  "G30",
			Value: "",
		}, {
			Cell:  "G31",
			Value: "",
		}, {
			Cell:  "G32",
			Value: "",
		}, {
			Cell:  "G33",
			Value: "",
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
			Value: "",
		}, {
			Cell:  "F67",
			Value: "",
		}, {
			Cell:  "F68",
			Value: "",
		},
		{
			Cell:  "F69",
			Value: "NO",
		}, {
			Cell:  "F70",
			Value: "",
		}, {
			Cell:  "C74",
			Value: "No",
		},
		{
			Cell:  "C93",
			Value: "",
		}, {
			Cell:  "G69",
			Value: "",
		},
		{
			Cell:  "E46",
			Value: "",
		},
	}
	return cells
}
