package quote

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/quote/internal"

	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/network"
	plc "gitlab.dev.wopta.it/goworkspace/policy"
	"gitlab.dev.wopta.it/goworkspace/product"
	"gitlab.dev.wopta.it/goworkspace/sellable"
)

func combinedQbeFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err        error
		reqPolicy  *models.Policy
		dbPolicy   models.Policy
		warrant    *models.Warrant
		inputCells []Cell
	)

	log.AddPrefix("CombinedQbeFx")
	defer func() {
		r.Body.Close()
		if err != nil {
			log.ErrorF("error: %s", err.Error())
		}
		log.Println("Handler end ---------------------------------------------")
		log.PopPrefix()
	}()
	log.Println("Handler start -----------------------------------------------")

	authToken, err := lib.GetAuthTokenFromIdToken(r.Header.Get("Authorization"))
	if err != nil {
		log.ErrorF("error getting authToken")
		return "", nil, err
	}
	log.Printf(
		"authToken - type: '%s' role: '%s' uid: '%s' email: '%s'",
		authToken.Type,
		authToken.Role,
		authToken.UserID,
		authToken.Email,
	)

	if err = json.NewDecoder(r.Body).Decode(&reqPolicy); err != nil {
		log.ErrorF("error decoding request body")
		return "", nil, err
	}

	if dbPolicy, err = plc.GetPolicy(reqPolicy.Uid); err != nil {
		log.ErrorF("error getting policy from DB")
		return "", nil, err
	}

	dbPolicy.Step = reqPolicy.Step
	dbPolicy.Assets = reqPolicy.Assets

	if err = sellable.CommercialCombined(&dbPolicy); err != nil {
		log.ErrorF("error on sellable")
		return "", nil, err
	}

	networkNode := network.GetNetworkNodeByUid(authToken.UserID)
	if networkNode != nil {
		warrant = networkNode.GetWarrant()
	}
	baseProduct := product.GetProductV2(dbPolicy.Name, dbPolicy.ProductVersion, dbPolicy.Channel, networkNode, warrant)

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
	mapCellPolicy(&dbPolicy, baseProduct, outCells, gsLink)

	if err = lib.SetFirestoreErr(lib.PolicyCollection, dbPolicy.Uid, dbPolicy); err != nil {
		log.ErrorF("error saving quote in policy")
		return "", nil, err
	}
	dbPolicy.BigquerySave()

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
	cells := make([]Cell, 0)
	cells = append(cells,
		Cell{Cell: totalNetYearlyCellValue},
		Cell{Cell: totalTaxAmountYearlyCellValue},
		Cell{Cell: totalGrossYearlyCellValue},
	)
	for _, column := range []string{grossYearlyValueColumn, netYearlyValueColumn,
		taxAmountYearlyValueColumn, taxValueColumn} {
		cells = append(cells,
			Cell{Cell: column + totalBuildingValueRow},
			Cell{Cell: column + totalStockValueRow},
			Cell{Cell: column + totalStockTemporaryIncreaseValueRow},
			Cell{Cell: column + totalTheftValueRow},
			Cell{Cell: column + totalRentalRiskValueRow},
			Cell{Cell: column + totalOtherContentValueRow},
			Cell{Cell: column + totalThirdPartyRecourseValueRow},
			Cell{Cell: column + totalFormulaValueRow},
			Cell{Cell: column + totalLossRentValueRow},
			Cell{Cell: column + totalThirdPartyLiabilityValueRow},
			Cell{Cell: column + totalWorkEmployersLiabilityValueRow},
			Cell{Cell: column + totalProductLiabilityValueRow},
			Cell{Cell: column + totalProductWithdrawalValueRow},
			Cell{Cell: column + totalManagementOrganizationValueRow},
			Cell{Cell: column + totalCyberValueRow},
		)
	}
	return cells
}

func mapCellByColumnAndSection(column, section string, priceGroup map[string]models.Price, cell Cell) {
	var hasError bool
	rawValue := cell.Value.(string)
	parsedValue, err := parseCellValue(cell.Value)
	if err != nil {
		hasError = true
		log.WarningF("error parsing value: %s", err.Error())
	}

	switch column {
	case grossYearlyValueColumn:
		if entry, ok := priceGroup[section]; ok && !hasError {
			entry.Gross = parsedValue
			priceGroup[section] = entry
		}
	case netYearlyValueColumn:
		if entry, ok := priceGroup[section]; ok && !hasError {
			entry.Net = parsedValue
			priceGroup[section] = entry
		}
	case taxAmountYearlyValueColumn:
		if entry, ok := priceGroup[section]; ok && !hasError {
			entry.Tax = parsedValue
			priceGroup[section] = entry
		}
	case taxValueColumn:
		return
	}

	if hasError {
		if entry, ok := priceGroup[section]; ok {
			entry.Description = rawValue
			priceGroup[section] = entry
		}
	}
}

func mapCellsToPriceGroup(cells []Cell) []models.Price {
	priceGroup := make([]models.Price, 0)
	priceGroupMap := make(map[string]models.Price)

	for key, value := range totalBySectionMap {
		priceGroupMap[key] = models.Price{Name: value}
	}

	for _, cell := range cells {
		cellColumn := cell.Cell[0:1]
		cellRow := cell.Cell[1:]
		mapCellByColumnAndSection(cellColumn, cellRow, priceGroupMap, cell)
	}

	for _, key := range totalBySectionOrder {
		priceGroup = append(priceGroup, priceGroupMap[key])
	}

	return priceGroup
}

func mapCellPolicy(policy *models.Policy, baseProduct *models.Product, cells []Cell, gsLink string) {
	log.AddPrefix("Commercial")
	defer log.PopPrefix()
	var (
		hasQuoteError bool
		quoteAtt      = models.Attachment{
			Name:      "QUOTAZIONE",
			FileName:  "Quotazione Excel.xlsx",
			MimeType:  "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
			Link:      gsLink,
			IsPrivate: true,
			Section:   "other",
			Note:      "",
		}
	)

	policy.OffersPrices = map[string]map[string]*models.Price{
		"default": {
			"yearly":  &models.Price{},
			"monthly": &models.Price{},
		},
	}

	policy.ReservedInfo.ReservedReasons = slices.DeleteFunc(policy.ReservedInfo.ReservedReasons, hasQuoteErrorFn)

	policy.PriceGroup = mapCellsToPriceGroup(cells)

	policyCells := slices.DeleteFunc(cells, func(c Cell) bool {
		return !slices.Contains([]string{totalNetYearlyCellValue, totalTaxAmountYearlyCellValue, totalGrossYearlyCellValue}, c.Cell)
	})

	for _, cell := range policyCells {
		parsedValue, err := parseCellValue(cell.Value)
		if err != nil {
			log.ErrorF("error parsing value: %s", err.Error())
			continue
		}
		switch cell.Cell {
		case totalNetYearlyCellValue:
			policy.OffersPrices["default"]["yearly"].Net = parsedValue
			policy.PriceNett = parsedValue
		case totalTaxAmountYearlyCellValue:
			policy.OffersPrices["default"]["yearly"].Tax = parsedValue
			policy.TaxAmount = parsedValue
		case totalGrossYearlyCellValue:
			policy.OffersPrices["default"]["yearly"].Gross = parsedValue
			policy.PriceGross = parsedValue
		}
	}

	log.Println("apply consultacy price")

	internal.AddConsultacyPrice(policy, baseProduct)

	if hasQuoteError {
		reserved := models.ReservedData{
			Id:          quoteErrorReservedDataId,
			Name:        "quote",
			Description: "Quotazione non effettuata",
		}
		policy.ReservedInfo.ReservedReasons = append(policy.ReservedInfo.ReservedReasons, reserved)
	}

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

func parseCellValue(value interface{}) (float64, error) {
	return strconv.ParseFloat(
		strings.Trim(strings.Replace(strings.Replace(value.(string), ".", "", -1), ",", ".", -1), " "), 64)
}

func hasQuoteErrorFn(r models.ReservedData) bool {
	return r.Id == quoteErrorReservedDataId
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

// Output Cells Column Indeces
const (
	grossYearlyValueColumn     = "L"
	netYearlyValueColumn       = "N"
	taxAmountYearlyValueColumn = "P"
	taxValueColumn             = "J"
)

// Output Cells Row Indeces
const (
	totalBuildingValueRow               = "81"
	totalStockValueRow                  = "82"
	totalStockTemporaryIncreaseValueRow = "83"
	totalTheftValueRow                  = "84"
	totalRentalRiskValueRow             = "85"
	totalOtherContentValueRow           = "86"
	totalThirdPartyRecourseValueRow     = "87"
	totalFormulaValueRow                = "88"
	totalLossRentValueRow               = "89"
	totalThirdPartyLiabilityValueRow    = "90"
	totalWorkEmployersLiabilityValueRow = "91"
	totalProductLiabilityValueRow       = "92"
	totalProductWithdrawalValueRow      = "93"
	totalManagementOrganizationValueRow = "94"
	totalCyberValueRow                  = "95"
)

// Output Cells Value Mapping
const (
	totalNetYearlyCellValue       = "C96"
	totalTaxAmountYearlyCellValue = "C97"
	totalGrossYearlyCellValue     = "C98"
)

// Price Group Section Names
const (
	totalBuildingPriceGroupTitle               = "Fabbricato"
	totalStockPriceGroupTitle                  = "Contenuto (Merci e Macchinari)"
	totalStockTemporaryIncreasePriceGroupTitle = "Merci (aumento temporaneo)"
	totalTheftPriceGroupTitle                  = "Furto, rapina, estorsione (in aumento)"
	totalRentalRiskPriceGroupTitle             = "Rischio locativo (in aumento)"
	totalOtherContentPriceGroupTitle           = "Altre garanzie su Contenuto"
	totalThirdPartyRecoursePriceGroupTitle     = "Ricorso terzi (in aumento)"
	totalFormulaPriceGroupTitle                = "Danni indiretti"
	totalLossRentPriceGroupTitle               = "Perdita Pigioni"
	totalThirdPartyLiabilityPriceGroupTitle    = "Responsabilità civile terzi"
	totalWorkEmployersLiabilityPriceGroupTitle = "Responsabilità civile prestatori lavoro"
	totalProductLiabilityPriceGroupTitle       = "Responsabilità civile prodotti"
	totalProductWithdrawalPriceGroupTitle      = "Ritiro Prodotti"
	totalManagementOrganizationPriceGroupTitle = "Resp. Amministratori Sindaci Dirigenti (D&O)"
	totalCyberPriceGroupTitle                  = "Cyber"
)

const (
	dateFormat               = "02/01/2006"
	quoteErrorReservedDataId = 10
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

var totalBySectionMap = map[string]string{
	totalBuildingValueRow:               totalBuildingPriceGroupTitle,
	totalStockValueRow:                  totalStockPriceGroupTitle,
	totalStockTemporaryIncreaseValueRow: totalStockTemporaryIncreasePriceGroupTitle,
	totalTheftValueRow:                  totalTheftPriceGroupTitle,
	totalRentalRiskValueRow:             totalRentalRiskPriceGroupTitle,
	totalOtherContentValueRow:           totalOtherContentPriceGroupTitle,
	totalThirdPartyRecourseValueRow:     totalThirdPartyRecoursePriceGroupTitle,
	totalFormulaValueRow:                totalFormulaPriceGroupTitle,
	totalLossRentValueRow:               totalLossRentPriceGroupTitle,
	totalThirdPartyLiabilityValueRow:    totalThirdPartyLiabilityPriceGroupTitle,
	totalWorkEmployersLiabilityValueRow: totalWorkEmployersLiabilityPriceGroupTitle,
	totalProductLiabilityValueRow:       totalProductLiabilityPriceGroupTitle,
	totalProductWithdrawalValueRow:      totalProductWithdrawalPriceGroupTitle,
	totalManagementOrganizationValueRow: totalManagementOrganizationPriceGroupTitle,
	totalCyberValueRow:                  totalCyberPriceGroupTitle,
}

var totalBySectionOrder = []string{
	totalBuildingValueRow,
	totalStockValueRow,
	totalStockTemporaryIncreaseValueRow,
	totalTheftValueRow,
	totalRentalRiskValueRow,
	totalOtherContentValueRow,
	totalThirdPartyRecourseValueRow,
	totalFormulaValueRow,
	totalLossRentValueRow,
	totalThirdPartyLiabilityValueRow,
	totalWorkEmployersLiabilityValueRow,
	totalProductLiabilityValueRow,
	totalProductWithdrawalValueRow,
	totalManagementOrganizationValueRow,
	totalCyberValueRow,
}
