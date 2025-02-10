package quote

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	plc "github.com/wopta/goworkspace/policy"
	"github.com/wopta/goworkspace/sellable"
)

func CombinedQbeFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err        error
		reqPolicy  *models.Policy
		dbPolicy   models.Policy
		inputCells []Cell
	)

	log.SetPrefix("[CombinedQbeFx] ")
	defer func() {
		r.Body.Close()
		if err != nil {
			log.Printf("error: %s", err.Error())
		}
		log.Println("Handler end ---------------------------------------------")
		log.SetPrefix("")
	}()
	log.Println("Handler start -----------------------------------------------")

	if err = json.NewDecoder(r.Body).Decode(&reqPolicy); err != nil {
		log.Println("error decoding request body")
		return "", nil, err
	}

	if dbPolicy, err = plc.GetPolicy(reqPolicy.Uid, ""); err != nil {
		log.Println("error getting policy from DB")
		return "", nil, err
	}

	dbPolicy.Assets = reqPolicy.Assets

	if err = sellable.CommercialCombined(&dbPolicy); err != nil {
		log.Println("error on sellable")
		return "", nil, err
	}

	inputCells = append(inputCells, setInputCell(&dbPolicy)...)
	qs := QuoteSpreadsheet{
		Id:                 "1tn0Jqce-r_JKdecExFOFVEJdGUaPYdGo31A9FOgvt-Y",
		DestinationSheetId: "1tMi7NYFZu7AnV4WkVrD0yzy1Dt3d-wVs0iZwlOcxLrg",
		InputCells:         inputCells,
		OutputCells:        setOutputCell(),
		InitCells:          resetCells(),
		SheetName:          "Input dati Polizza",
		ExportedSheetName:  "Export",
		ExportFilePrefix:   fmt.Sprintf("quote_%s_%s", dbPolicy.Name, dbPolicy.Uid),
	}
	outCells, gsLink := qs.Spreadsheets()
	mapCellPolicy(&dbPolicy, outCells, gsLink)

	if err = lib.SetFirestoreErr(lib.PolicyCollection, dbPolicy.Uid, dbPolicy); err != nil {
		log.Println("error saving quote in policy")
		return "", nil, err
	}
	dbPolicy.BigquerySave("")

	policyJson, err := dbPolicy.Marshal()

	return string(policyJson), dbPolicy, err
}

func resetCells() []Cell {
	headingCellList := []string{producerValueCell, enterpriseNameValueCell, startDateValueCell, endDateValueCell, vatCodeValueCell}
	buildingsCellList := make([]string, 0)
	for _, column := range buildingMap {
		buildingsCellList = append(buildingsCellList,
			column+naicsCategoryValueRow, column+naicsDetailValueRow, column+naicsValueRow,
			column+postalCodeValueRow, column+provinceValueRow, column+cityValueRow, column+addressValueRow,
			column+buildingMaterialValueRow, column+sandwichPanelValueRow, column+alarmValueRow, column+sprinklerValueRow,
			column+buildingValueRow, column+rentalRiskValueRow, column+machineryValueRow, column+stockValueRow, column+stockTemporaryIncreaseValueRow,
		)
	}
	enterpriseCellList := []string{
		employeeNumberValueCell, remunerationValueCell, revenueValueCell, revenueUsaCanValueCell,
		stockTemporaryIncreaseDurationCell, thirdPartyRecourseValueCell, electricalPhenomenonValueCell,
		refrigerationStockValueCell, machineryBreakdownValueCell, electronicEquipmentValueCell, theftValueCell,
		dailyAllowanceValueCell, increasedCostValueCell, lossRentValueCell, stockTemporaryIncreaseDateCell,
		dailyAllowanceDurationCell, thirdPartyLiabilityWorkProvidersValueCell, managementOrganizationTotalAssetCell,
		managementOrganizationOwnCapitalCell, thirdPartyLiabilityWorkProvidersRetroactiveDateCell,
		productLiabilityRetroactiveDateCell, productLiabilityRetroactiveDateUsaCanCell, managementOrganizationDateCell,
		productWithdrawalDateCell,
	}
	discountCellList := []string{discountGoodsValueCell, discountTheftValueCell, discountLiabilityValueCell}

	toEmptyCells := make([]string, 0)
	toEmptyCells = append(toEmptyCells, headingCellList...)
	toEmptyCells = append(toEmptyCells, buildingsCellList...)
	toEmptyCells = append(toEmptyCells, enterpriseCellList...)
	toEmptyCells = append(toEmptyCells, discountCellList...)

	initializedCells := []Cell{
		{Cell: firstRateMergerValueCell, Value: yesValue},
		{Cell: paymentSplitValueCell, Value: paymentSplitYearlyValue},
		{Cell: bondValueCell, Value: noValue},
		{Cell: formulaValueCell, Value: formulaExcludedValue},
		{Cell: productLiabilityValueCell, Value: noValue},
		{Cell: managementOrganizationValueCell, Value: noValue},
		{Cell: productWithdrawalChoiceCell, Value: noValue},
		{Cell: cyberValueCell, Value: noValue},
	}

	for _, cell := range toEmptyCells {
		initializedCells = append(initializedCells, Cell{cell, emptyValue})
	}

	return initializedCells
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
		FileName:  "Quotazione Excel.xlsx",
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
				Id:          10,
				Name:        "quote",
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

	if policy.Attachments == nil {
		policy.Attachments = new([]models.Attachment)
	}

	quoteAttIdx := slices.IndexFunc(*policy.Attachments, func(a models.Attachment) bool {
		return a.Name == quoteAtt.Name
	})

	if quoteAttIdx == -1 {
		*policy.Attachments = append(*policy.Attachments, quoteAtt)
	} else {
		(*policy.Attachments)[quoteAttIdx].Link = gsLink
	}
}

func setInputCell(policy *models.Policy) []Cell {
	var inputCells []Cell

	assEnterprise := getAssetByType(policy, models.AssetTypeEnterprise)
	assBuildings := getAssetByType(policy, models.AssetTypeBuilding)

	inputCells = append(inputCells, Cell{Cell: startDateValueCell, Value: policy.StartDate.Format(dateFormat)})
	inputCells = append(inputCells, Cell{Cell: paymentSplitValueCell, Value: paymentSplitMap[policy.PaymentSplit]})
	inputCells = append(inputCells, Cell{Cell: endDateValueCell, Value: policy.EndDate.Format(dateFormat)})

	inputCells = append(inputCells, setEnterpriseCell(assEnterprise[0])...)

	for i, build := range assBuildings {
		buildingColumn := buildingMap[i]
		inputCells = append(inputCells, Cell{Cell: buildingColumn + naicsCategoryValueRow, Value: build.Building.NaicsCategory})
		inputCells = append(inputCells, Cell{Cell: buildingColumn + naicsDetailValueRow, Value: build.Building.NaicsDetail})
		inputCells = append(inputCells, Cell{Cell: buildingColumn + naicsValueRow, Value: build.Building.Naics})
		inputCells = append(inputCells, Cell{Cell: buildingColumn + postalCodeValueRow, Value: build.Building.BuildingAddress.PostalCode})
		inputCells = append(inputCells, Cell{Cell: buildingColumn + provinceValueRow, Value: build.Building.BuildingAddress.Locality})
		inputCells = append(inputCells, Cell{Cell: buildingColumn + cityValueRow, Value: build.Building.BuildingAddress.City})
		inputCells = append(inputCells, Cell{Cell: buildingColumn + addressValueRow, Value: build.Building.BuildingAddress.StreetName})
		inputCells = append(inputCells, Cell{Cell: buildingColumn + buildingMaterialValueRow, Value: build.Building.BuildingMaterial})
		inputCells = append(inputCells, Cell{Cell: buildingColumn + sandwichPanelValueRow, Value: booleanMap[build.Building.HasSandwichPanel]})
		inputCells = append(inputCells, Cell{Cell: buildingColumn + alarmValueRow, Value: booleanMap[build.Building.HasAlarm]})
		inputCells = append(inputCells, Cell{Cell: buildingColumn + sprinklerValueRow, Value: booleanMap[build.Building.HasSprinkler]})
		for _, bg := range build.Guarantees {
			inputCells = append(inputCells, getBuildingGuaranteCellsBySlug(bg, buildingColumn)...)
		}
	}

	for _, c := range inputCells {
		log.Printf("Cell: '%s' - Value: '%+v'", c.Cell, c.Value)
	}

	return inputCells
}

func setEnterpriseCell(asset models.Asset) []Cell {
	var inputCells []Cell
	inputCells = append(inputCells, Cell{Cell: vatCodeValueCell, Value: asset.Enterprise.VatCode})
	inputCells = append(inputCells, Cell{Cell: enterpriseNameValueCell, Value: asset.Enterprise.Name})
	inputCells = append(inputCells, Cell{Cell: employeeNumberValueCell, Value: asset.Enterprise.Employer})
	inputCells = append(inputCells, Cell{Cell: remunerationValueCell, Value: asset.Enterprise.WorkEmployersRemuneration})
	inputCells = append(inputCells, Cell{Cell: revenueValueCell, Value: asset.Enterprise.TotalBilled})
	inputCells = append(inputCells, Cell{Cell: revenueUsaCanValueCell, Value: asset.Enterprise.NorthAmericanMarket})
	for _, eg := range asset.Guarantees {
		inputCells = append(inputCells, getEnterpriseGuaranteCellsBySlug(eg)...)
	}
	return inputCells
}

func getAssetByType(policy *models.Policy, assetType string) []models.Asset {
	var assets []models.Asset
	for _, asset := range policy.Assets {
		if asset.Type == assetType {
			assets = append(assets, asset)
		}
	}
	return assets
}

func getEnterpriseGuaranteCellsBySlug(guarante models.Guarante) []Cell {
	var cells []Cell
	switch guarante.Slug {
	case electricalPhenomenonGuaranteeSlug:
		cells = []Cell{{
			Cell:  electricalPhenomenonValueCell,
			Value: guarante.Value.SumInsuredLimitOfIndemnity,
		}}
	case refrigerationStockGuaranteeSlug:
		cells = []Cell{{
			Cell:  refrigerationStockValueCell,
			Value: guarante.Value.SumInsuredLimitOfIndemnity,
		}}
	case machineryBreakdownGuaranteeSlug:
		cells = []Cell{{
			Cell:  machineryBreakdownValueCell,
			Value: guarante.Value.SumInsuredLimitOfIndemnity,
		}}
	case electronicEquipmentGuaranteeSlug:
		cells = []Cell{{
			Cell:  electronicEquipmentValueCell,
			Value: guarante.Value.SumInsuredLimitOfIndemnity,
		}}
	case theftGuaranteeSlug:
		cells = []Cell{{
			Cell:  theftValueCell,
			Value: guarante.Value.SumInsuredLimitOfIndemnity,
		}, {
			Cell:  theftDiscountCell,
			Value: int(guarante.Value.Discount),
		}}
	case thirdPartyRecourseGuaranteeSlug:
		cells = []Cell{{
			Cell:  thirdPartyRecourseValueCell,
			Value: guarante.Value.SumInsuredLimitOfIndemnity,
		}}
	case thirdPartyLiabilityWorkProvidersGuaranteeSlug:
		cells = []Cell{{
			Cell:  thirdPartyLiabilityWorkProvidersValueCell,
			Value: guarante.Value.SumInsuredLimitOfIndemnity,
		}, {
			Cell:  thirdPartyLiabilityWorkProvidersRetroactiveDateCell,
			Value: guarante.Value.RetroactiveDate.Format(dateFormat),
		}, {
			Cell:  thirdPartyLiabilityWorkProvidersDiscountCell,
			Value: int(guarante.Value.Discount),
		}}
	case productLiabilityGuaranteeSlug:
		cells = []Cell{{
			Cell:  productLiabilityValueCell,
			Value: guarante.Value.SumInsuredLimitOfIndemnity,
		}, {
			Cell:  productLiabilityRetroactiveDateCell,
			Value: guarante.Value.RetroactiveDate.Format(dateFormat),
		}, {
			Cell:  productLiabilityRetroactiveDateUsaCanCell,
			Value: guarante.Value.RetroactiveUsaCanDate.Format(dateFormat),
		}}
	case managementOrganizationGuaranteeSlug:
		cells = []Cell{{
			Cell:  managementOrganizationValueCell,
			Value: guarante.Value.SumInsuredLimitOfIndemnity,
		}, {
			Cell:  managementOrganizationTotalAssetCell,
			Value: guarante.Value.LimitOfIndemnity,
		}, {
			Cell:  managementOrganizationOwnCapitalCell,
			Value: guarante.Value.SumInsured,
		}, {
			Cell:  managementOrganizationDateCell,
			Value: guarante.Value.StartDate.Format(dateFormat),
		}}
	case cyberGuranteeSlug:
		cells = []Cell{{
			Cell:  cyberValueCell,
			Value: guarante.Value.SumInsuredLimitOfIndemnity,
		}}
	case dailyAllowanceGuaranteeSlug:
		cells = []Cell{{
			Cell:  dailyAllowanceValueCell,
			Value: guarante.Value.SumInsuredLimitOfIndemnity,
		}, {
			Cell:  formulaValueCell,
			Value: formulaDailyAllowanceValue,
		}, {
			Cell:  dailyAllowanceDurationCell,
			Value: guarante.Value.Duration.Day,
		}}
	case additionalCompensationGuaranteeSlug:
		cells = []Cell{{
			Cell:  formulaValueCell,
			Value: formulaAdditionalCompensationValue,
		}}
	case increasedCostGuaranteeSlug:
		cells = []Cell{{
			Cell:  increasedCostValueCell,
			Value: guarante.Value.SumInsuredLimitOfIndemnity,
		}, {
			Cell:  formulaValueCell,
			Value: formulaIncreasedCostValue,
		}}
	case lossRentGuaranteeSlug:
		cells = []Cell{{
			Cell:  lossRentValueCell,
			Value: guarante.Value.SumInsuredLimitOfIndemnity,
		}}
	case productWithdrawalGuaranteeSlug:
		cells = []Cell{{
			Cell:  productWithdrawalChoiceCell,
			Value: yesValue,
		}, {
			Cell:  productWithdrawalDateCell,
			Value: guarante.Value.StartDate.Format(dateFormat),
		}}
	}
	return cells
}

func getBuildingGuaranteCellsBySlug(guarante models.Guarante, buildingColumn string) []Cell {
	var cells []Cell
	switch guarante.Slug {
	case buildingGuaranteeSlug:
		cells = []Cell{{
			Cell:  buildingColumn + buildingValueRow,
			Value: guarante.Value.SumInsuredLimitOfIndemnity,
		}}
	case rentalRiskGuaranteeSlug:
		cells = []Cell{{
			Cell:  buildingColumn + rentalRiskValueRow,
			Value: guarante.Value.SumInsuredLimitOfIndemnity,
		}}
	case machineryGuaranteeSlug:
		cells = []Cell{{
			Cell:  buildingColumn + machineryValueRow,
			Value: guarante.Value.SumInsuredLimitOfIndemnity,
		}}
	case stockGuaranteeSlug:
		cells = []Cell{{
			Cell:  buildingColumn + stockValueRow,
			Value: guarante.Value.SumInsuredLimitOfIndemnity,
		}}
	case stockTemporaryIncreaseGuaranteeSlug:
		cells = []Cell{{
			Cell:  buildingColumn + stockTemporaryIncreaseValueRow,
			Value: guarante.Value.SumInsuredLimitOfIndemnity,
		}, {
			Cell:  stockTemporaryIncreaseDateCell,
			Value: guarante.Value.StartDateString,
		}, {
			Cell:  stockTemporaryIncreaseDurationCell,
			Value: guarante.Value.Duration.Day,
		}}
	}
	return cells
}

// Enterprise Guarantees Slugs
const (
	electricalPhenomenonGuaranteeSlug             string = "electrical-phenomenon"
	refrigerationStockGuaranteeSlug               string = "refrigeration-stock"
	machineryBreakdownGuaranteeSlug               string = "machinery-breakdown"
	electronicEquipmentGuaranteeSlug              string = "electronic-equipment"
	theftGuaranteeSlug                            string = "theft"
	thirdPartyRecourseGuaranteeSlug               string = "third-party-recourse"
	thirdPartyLiabilityWorkProvidersGuaranteeSlug string = "third-party-liability-work-providers"
	productLiabilityGuaranteeSlug                 string = "product-liability"
	managementOrganizationGuaranteeSlug           string = "management-organization"
	cyberGuranteeSlug                             string = "cyber"
	dailyAllowanceGuaranteeSlug                   string = "daily-allowance"
	additionalCompensationGuaranteeSlug           string = "additional-compensation"
	increasedCostGuaranteeSlug                    string = "increased-cost"
	lossRentGuaranteeSlug                         string = "loss-rent"
	productWithdrawalGuaranteeSlug                string = "product-withdrawal"
)

// Enterprise Guarantees Cell Mapping
const (
	electricalPhenomenonValueCell                       = "C48"
	refrigerationStockValueCell                         = "C49"
	machineryBreakdownValueCell                         = "C50"
	electronicEquipmentValueCell                        = "C51"
	theftValueCell                                      = "C52"
	theftDiscountCell                                   = "G81"
	thirdPartyRecourseValueCell                         = "C47"
	thirdPartyLiabilityWorkProvidersValueCell           = "C66"
	thirdPartyLiabilityWorkProvidersRetroactiveDateCell = "F66"
	thirdPartyLiabilityWorkProvidersDiscountCell        = "G82"
	productLiabilityValueCell                           = "C67"
	productLiabilityRetroactiveDateCell                 = "F67"
	productLiabilityRetroactiveDateUsaCanCell           = "F68"
	managementOrganizationValueCell                     = "C68"
	managementOrganizationTotalAssetCell                = "C69"
	managementOrganizationOwnCapitalCell                = "C70"
	managementOrganizationDateCell                      = "F70"
	cyberValueCell                                      = "C74"
	dailyAllowanceValueCell                             = "C58"
	dailyAllowanceDurationCell                          = "E58"
	increasedCostValueCell                              = "C59"
	formulaValueCell                                    = "C57"
	lossRentValueCell                                   = "C61"
	productWithdrawalChoiceCell                         = "F69"
	productWithdrawalDateCell                           = "G69"
)

// Building Guarantees Slugs
const (
	buildingGuaranteeSlug               string = "building"
	rentalRiskGuaranteeSlug             string = "rental-risk"
	machineryGuaranteeSlug              string = "machinery"
	stockGuaranteeSlug                  string = "stock"
	stockTemporaryIncreaseGuaranteeSlug string = "stock-temporary-increase"
)

// Building Number Column Indeces
const (
	building1ValueColumn = "C"
	building2ValueColumn = "D"
	building3ValueColumn = "E"
	building4ValueColumn = "F"
	building5ValueColumn = "G"
)

// Building Guarantees Row Indeces
const (
	naicsCategoryValueRow          = "19"
	naicsDetailValueRow            = "20"
	naicsValueRow                  = "21"
	postalCodeValueRow             = "29"
	provinceValueRow               = "30"
	cityValueRow                   = "31"
	addressValueRow                = "32"
	buildingMaterialValueRow       = "33"
	sandwichPanelValueRow          = "34"
	alarmValueRow                  = "35"
	sprinklerValueRow              = "36"
	buildingValueRow               = "41"
	rentalRiskValueRow             = "42"
	machineryValueRow              = "43"
	stockValueRow                  = "44"
	stockTemporaryIncreaseValueRow = "45"
)

// Global Building Guarantees Cell mapping
const (
	stockTemporaryIncreaseDateCell     = "E46"
	stockTemporaryIncreaseDurationCell = "C46"
)

// Global Policy Cell Mapping
const (
	producerValueCell          = "C4"
	enterpriseNameValueCell    = "C5"
	vatCodeValueCell           = "E6"
	startDateValueCell         = "C10"
	endDateValueCell           = "C11"
	paymentSplitValueCell      = "C16"
	employeeNumberValueCell    = "C24"
	remunerationValueCell      = "C25"
	revenueValueCell           = "C26"
	revenueUsaCanValueCell     = "C27"
	discountGoodsValueCell     = "G80"
	discountTheftValueCell     = "G81"
	discountLiabilityValueCell = "G82"
	firstRateMergerValueCell   = "C14"
	bondValueCell              = "G11"
)

// Standard Values
const (
	formulaDailyAllowanceValue         = "Diaria Giornaliera"
	formulaAdditionalCompensationValue = "Indennità Aggiuntiva (+10%)"
	formulaIncreasedCostValue          = "Maggiori Costi"
	formulaExcludedValue               = "Esclusa"
	yesValue                           = "SI"
	noValue                            = "NO"
	emptyValue                         = ""
	paymentSplitYearlyValue            = "Annuale"
	paymentSplitSemestralValue         = "Semestrale"
)

const (
	dateFormat = "02/01/2006"
)

var buildingMap = map[int]string{
	0: building1ValueColumn,
	1: building2ValueColumn,
	2: building3ValueColumn,
	3: building4ValueColumn,
	4: building5ValueColumn,
}

var booleanMap = map[bool]string{
	true:  yesValue,
	false: noValue,
}

var paymentSplitMap = map[string]string{
	string(models.PaySplitYearly):    paymentSplitYearlyValue,
	string(models.PaySplitSemestral): paymentSplitSemestralValue,
}
