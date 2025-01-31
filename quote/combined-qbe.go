package quote

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"slices"
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
		Id:                 "1tn0Jqce-r_JKdecExFOFVEJdGUaPYdGo31A9FOgvt-Y",
		DestinationSheetId: "1tMi7NYFZu7AnV4WkVrD0yzy1Dt3d-wVs0iZwlOcxLrg",
		InputCells:         inputCells,
		OutputCells:        setOutputCell(),
		InitCells:          resetCells(),
		SheetName:          "Input dati Polizza",
		ExportedSheetName:  "Export",
		ExportFilePrefix:   fmt.Sprintf("quote_%s_%s", policy.Name, policy.Uid),
	}
	outCells, gsLink := qs.Spreadsheets()
	mapCellPolicy(policy, outCells, gsLink)

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
func mapCellPolicy(policy *models.Policy, cells []Cell, gsLink string) {
	var priceGroup []models.Price

	var quoteAtt = models.Attachment{
		Name:      "QUOTAZIONE",
		FileName:  "Quotazione Excel",
		MimeType:  "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		Link:      gsLink,
		IsPrivate: true,
		Section:   "other",
		Note:      "",
	}

	policy.OffersPrices = map[string]map[string]*models.Price{
		"default": {
			"yearly":  &models.Price{},
			"monthly": &models.Price{},
		},
	}

	for _, cell := range cells {
		v := cell.Value.(string)
		if strings.HasPrefix(v, "Errore") {
			reserved := models.ReservedData{
				Id: 10,
				Name: "quote",
				Description: "Quotazione non effettuata",
			}
			if !slices.ContainsFunc(policy.ReservedInfo.ReservedReasons, func(r models.ReservedData) bool {
				return r.Id == 10
			}) {
				policy.ReservedInfo.ReservedReasons = append(policy.ReservedInfo.ReservedReasons, reserved)
			}
		}

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

	if len(*policy.Attachments) == 0 {
		*policy.Attachments = append(*policy.Attachments, quoteAtt)
	} else {
		for i := 0; i < len(*policy.Attachments); i++ {
			if (*policy.Attachments)[i].Name == quoteAtt.Name {
				(*policy.Attachments)[i].Link = gsLink
			}
		}
	}

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
			inputCells = append(inputCells, Cell{Cell: col[i] + "29", Value: build.Building.BuildingAddress.PostalCode})
			inputCells = append(inputCells, Cell{Cell: col[i] + "30", Value: build.Building.BuildingAddress.Locality})
			inputCells = append(inputCells, Cell{Cell: col[i] + "31", Value: build.Building.BuildingAddress.City})
			inputCells = append(inputCells, Cell{Cell: col[i] + "32", Value: build.Building.BuildingAddress.StreetName})
			inputCells = append(inputCells, Cell{Cell: col[i] + "19", Value: build.Building.NaicsCategory})
			inputCells = append(inputCells, Cell{Cell: col[i] + "20", Value: build.Building.NaicsDetail})
			inputCells = append(inputCells, Cell{Cell: col[i] + "21", Value: build.Building.Naics})
			inputCells = append(inputCells, Cell{Cell: col[i] + "33", Value: build.Building.BuildingMaterial})
			if build.Building.HasSandwichPanel {
				inputCells = append(inputCells, Cell{Cell: col[i] + "34", Value: "SI"})
			} else {
				inputCells = append(inputCells, Cell{Cell: col[i] + "34", Value: "NO"})
			}
			if build.Building.HasAlarm {
				inputCells = append(inputCells, Cell{Cell: col[i] + "35", Value: "SI"})
			} else {
				inputCells = append(inputCells, Cell{Cell: col[i] + "35", Value: "NO"})
			}
			if build.Building.HasSprinkler {
				inputCells = append(inputCells, Cell{Cell: col[i] + "36", Value: "SI"})
			} else {
				inputCells = append(inputCells, Cell{Cell: col[i] + "36", Value: "NO"})
			}
			inputCells = append(inputCells, getBuildingGuaranteCellsBySlug(bg, i)...)
		}

	}

	for _, c := range inputCells {
		log.Printf("Cell: '%s' - Value: '%+v'", c.Cell, c.Value)
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
	case "excluded-formula":
		cells = []Cell{
			{
				Cell:  "C57",
				Value: "Esclusa",
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

func resetCells() []Cell {
	headingCellList := []string{"C4", "C5", "C6", "C7", "C8", "C9", "C10", "C11", "C15", "D9", "E6", "G6", "G8", "G9", "G10"}
	buildingsCellList := []string{
		"C19", "C20", "C21", "C29", "C30", "C31", "C32", "C33", "C34", "C35", "C36", "C41", "C42", "C43", "C44", "C45",
		"D19", "D20", "D21", "D29", "D30", "D31", "D32", "D33", "D34", "D35", "D36", "D41", "D42", "D43", "D44", "D45",
		"E19", "E20", "E21", "E29", "E30", "E31", "E32", "E33", "E34", "E35", "E36", "E41", "E42", "E43", "E44", "E45",
		"F19", "F20", "F21", "F29", "F30", "F31", "F32", "F33", "F34", "F35", "F36", "F41", "F42", "F43", "F44", "F45",
		"G19", "G20", "G21", "G29", "G30", "G31", "G32", "G33", "G34", "G35", "G36", "G41", "G42", "G43", "G44", "G45",
	}
	enterpriseCellList := []string{
		"C24", "C25", "C26", "C27",
		"C46", "C47", "C48", "C49", "C50", "C51", "C52",
		"C58", "C59", "C61", "E58",
		"C66", "C69", "C70", "F66", "F67", "F68", "F70", "G69",
		"C93",
	}
	discountCellList := []string{"G80", "G81", "G82"}

	allCells := make([]string, 0)
	allCells = append(allCells, headingCellList...)
	allCells = append(allCells, buildingsCellList...)
	allCells = append(allCells, enterpriseCellList...)
	allCells = append(allCells, discountCellList...)

	initializedCells := []Cell{
		{Cell: "C14", Value: "SI"},
		{Cell: "C16", Value: "Annuale"},
		{Cell: "G11", Value: "NO"},
		{Cell: "C57", Value: "Esclusa"},
		{Cell: "C67", Value: "NO"},
		{Cell: "C68", Value: "NO"},
		{Cell: "F69", Value: "NO"},
		{Cell: "C74", Value: "NO"},
	}

	for _, cell := range allCells {
		initializedCells = append(initializedCells, Cell{cell, ""})
	}

	return initializedCells
}
