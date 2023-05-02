package companydata

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	lib "github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func Emit(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		policy models.Policy
	)

	e := json.Unmarshal([]byte(getpolicymock()), &policy)
	m := make(map[string]interface{})
	FieldNames(policy, m)
	for name, v := range m {
		fmt.Println(name)
		fmt.Println(v)
	}

	return "", nil, e
}
func getFieldValue(v interface{}, field string) string {
	r := reflect.ValueOf(v)

	f := reflect.Indirect(r).FieldByName(field)

	return f.String()
}
func GetStructFieldName(Struct interface{}, StructField ...interface{}) (fields map[int]string) {
	fields = make(map[int]string)
	s := reflect.ValueOf(Struct).Elem()

	for r := range StructField {
		f := reflect.ValueOf(StructField[r]).Elem()

		for i := 0; i < s.NumField(); i++ {
			valueField := s.Field(i)
			if valueField.Addr().Interface() == f.Addr().Interface() {
				fields[i] = s.Type().Field(i).Name
			}
		}
	}
	return fields
}
func FieldNames(Struct interface{}, m map[string]interface{}) {
	v := reflect.ValueOf((Struct))
	t := reflect.TypeOf((Struct))
	//element:=v.Elem()
	if t.Kind() == reflect.Struct {
		for i := 0; i < t.NumField(); i++ {
			fieldValue := v.Field(i)
			typeValue := t.Field(i)
			fmt.Println(typeValue.Name)
			fmt.Println(typeValue.Type)
			fmt.Println(fieldValue.Interface())
			if typeValue.Type.Kind() == reflect.Struct {
				FieldNames(fieldValue.Interface(), m)

			}
		}
	}
}
func collectFieldNames(t reflect.Type, m map[string]interface{}) {
	if t.Kind() == reflect.Ptr {

		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		m[sf.Name] = t.Field(i)
		if sf.Anonymous {
			collectFieldNames(sf.Type, m)
		}
	}
}
func identify(output map[string]interface{}) {
	fmt.Printf("%T", output)
	for _, b := range output {
		switch bb := b.(type) {
		case string:
			fmt.Println(bb)
			fmt.Println("This is a string")
		case float64:
			fmt.Println("this is a float")
		case bool:
			fmt.Println("this is a boolean")
		case []interface{}:
		// Access the values in the JSON object and place them in an Item

		default:
			return
		}
	}
}
func getPolicy() []models.Policy {

	q := lib.Firequeries{
		Queries: []lib.Firequery{{
			Field:      "companyEmit", //
			Operator:   "==",          //
			QueryValue: true,
		},
			{
				Field:      "companyEmitted", //
				Operator:   "==",             //
				QueryValue: false,
			},
		},
	}
	query, _ := q.FirestoreWherefields("policy")
	policies := models.PolicyToListData(query)
	return policies
}

func getpolicymock() string {

	return `{
    "startDate": "2023-03-23T21:16:35.982Z",
    "endDate": "2024-03-23T21:16:35.982Z",
    "offerName": "Completa",
    "status": "",
    "statusHistory": [],
    "transactions": [],
    "company": "global",
    "name": "pmi",
    "payment":"fabrik",
    "paymentType": "creditCard",
    "paymentSplit":"year",
    "coverageType": "",
    "voucher": "",
    "channel": "",
    "covenant": "",
    "price": 0,
    "priceNett": 1099.99,
    "priceGross": 1337.06,
    "contractor": {
        "name": "Beatrice",
        "surname": "Sala",
        "birthDate": "2005-03-01T23:00:00.000Z",
        "birthProvince": "MI",
        "birthCity": "Milano",
        "fiscalCode": "brblcu81h03f205q",
        "mail": "luca.barbieri@wopta.it",
        "phone": "+393668134257",
        "role": "",
        "work": "",
        "workType": "",
        "type": "",
        "cluster": "",
        "riskClass": "",
        "vatCode": "01319960199",
        "address": "Galleria del Corso",
        "postalCode": "20122",
        "city": "Milano",
        "cityCode": "MI",
        "streetNumber": "1",
        "location": {
            "lat": 45.4648158,
            "lng": 9.1949345
        },
        "consens": [
            {
                "key": 1,
                "title": "Privacy",
                "consens": "Il sottoscritto, letta e compresa l'informativa sul trattamento dei dati personali, ACCONSENTE al trattamento dei propri dati personali da parte di Wopta Assicurazioni per l’invio di comunicazioni e proposte commerciali e di marketing, incluso l’invio di newsletter e ricerche di mercato, attraverso strumenti automatizzati (sms, mms, e-mail, ecc.) e non (posta cartacea e telefono con operatore).",
                "answer": true
            }
        ]
    },
    "attachments": [
        {}
    ],
    "assets": [
        {
            "enterprise": {
                "name": "Beatrice",
                "type": "",
                "vatCode": "01319960199",
                "ateco": "10.71.10",
                "atecoDesc": "PRODUZIONE DI PRODOTTI DI PANETTERIA FRESCHI",
                "atecoMacro": "ATTIVITA MANIFATTURIERE",
                "atecoSub": "INDUSTRIE ALIMENTARI",
                "revenue": "270383",
                "employer": 2,
                "class": "",
                "sector": "PRODUZIONE",
                "address": "Galleria del Corso",
                "postalCode": "20122",
                "city": "Milano",
                "cityCode": "MI",
                "streetNumber": "1",
                "location": {
                    "lat": 45.4648158,
                    "lng": 9.1949345
                }
            },
            "guarantees": [
                {
                    "Assistance": "",
                    "CompanyCodec": "DP",
                    "CompanyName": "Tutela Legale",
                    "Config": null,
                    "DailyAllowance": "",
                    "Deductible": "",
                    "Description": "Difesa penale per reati di natura colposa o contravvenzionale, inclusi i casi di sicurezza aziendale da D. Lgs. 81/08 e D. Lgs. 106/09, D. Lgs. 193/07, D. Lgs. 152/06, D. Lgs. 101/18, D. Lgs. 231/01",
                    "Group": "LEGAL",
                    "IsBase": false,
                    "IsExtension": false,
                    "IsPremium": true,
                    "IsSellable": false,
                    "IsYour": false,
                    "IsYuor": false,
                    "LegalDefence": "basic",
                    "Name": "",
                    "Offer": null,
                    "Price": 0,
                    "PriceGross": 50.0035,
                    "PriceNett": 41.24,
                    "SelfInsurance": "",
                    "SelfInsuranceDesc": "",
                    "Slug": "legal-defence",
                    "SumInsuredLimitOfIndemnity": 0,
                    "Tax": 21.25,
                    "Taxes": null,
                    "Type": "company",
                    "TypeOfSumInsured": "",
                    "Value": null,
                    "translation": "Tutela legale"
                },
                {
                    "Assistance": "",
                    "CompanyCodec": "CO",
                    "CompanyName": "Resp. Civile verso Prestatori di Lavoro",
                    "Config": null,
                    "DailyAllowance": "",
                    "Deductible": "2600",
                    "Description": "Responsabilità per: la rivalsa INAIL per gli infortuni sul lavoro subiti dai prestatori di lavoro;  morte; e lesioni personali dalle quali sia derivata un'invalidità permanente ai sensi del codice civile, incluse le malattie professionali",
                    "Group": "RCT",
                    "IsBase": true,
                    "IsExtension": false,
                    "IsPremium": true,
                    "IsSellable": false,
                    "IsYour": false,
                    "IsYuor": true,
                    "LegalDefence": "",
                    "Name": "",
                    "Offer": null,
                    "Price": 0,
                    "PriceGross": 77.836575,
                    "PriceNett": 63.67,
                    "SelfInsurance": "0.10",
                    "SelfInsuranceDesc": "10% - minimo € 1500",
                    "Slug": "employers-liability",
                    "SumInsuredLimitOfIndemnity": 1000000,
                    "Tax": 22.25,
                    "Taxes": null,
                    "Type": "company",
                    "TypeOfSumInsured": "",
                    "Value": null,
                    "translation": "Responsabilità Civile Addetti"
                },
                {
                    "Assistance": "",
                    "CompanyCodec": "CT",
                    "CompanyName": "Responsabilità Civile Terzi",
                    "Config": null,
                    "DailyAllowance": "",
                    "Deductible": "1000",
                    "Description": "Danni involontariamente causati a terzi per danni a cose o persone di cui sia responsabile a termini di legge (R.C.T.). La garanzia include, ma non si limita a questi, i  danni: a veicoli di terzi e prestatori di lavoro; a cose in consegna e custodia; a cose nell'ambito di esecuzione dei lavori; a cose di terzi sollevate, caricate, scaricate, movimentate, trasportate o rimorchiate; a mezzi di trasporto sotto carico e scarico; da interruzione o sospensione di attività di terzi; da smercio; da committenza autoveicoli; da responsabilità civile personale addetti; da attività di commercio ambulante; da lavori presso terzi",
                    "Group": "RCT",
                    "IsBase": true,
                    "IsExtension": false,
                    "IsPremium": true,
                    "IsSellable": false,
                    "IsYour": false,
                    "IsYuor": true,
                    "LegalDefence": "",
                    "Name": "",
                    "Offer": null,
                    "Price": 0,
                    "PriceGross": 102.482175,
                    "PriceNett": 83.83,
                    "SelfInsurance": "",
                    "SelfInsuranceDesc": "",
                    "Slug": "third-party-liability",
                    "SumInsuredLimitOfIndemnity": 1000000,
                    "Tax": 22.25,
                    "Taxes": null,
                    "Type": "company",
                    "TypeOfSumInsured": "",
                    "Value": null,
                    "translation": "Responsabilità Civile Terzi"
                }
            ]
        },
        {
            "building": {
                "name": "Wopta",
                "type": "",
                "isAllarm": true,
                "isHolder": true,
                "buildingType": "industriale",
                "buildingMaterial": "masonry",
                "buildingYear": "before1972",
                "squareMeters": 200,
                "floor": "ground_floor",
                "construction": "",
                "address": "Galleria del Corso",
                "streetNumber": "1",
                "postalCode": "20122",
                "city": "Milano",
                "cityCode": "MI",
                "location": {
                    "lat": 45.4648158,
                    "lng": 9.1949345
                }
            },
            "guarantees": [
                {
                    "Assistance": "yes",
                    "CompanyCodec": "AS",
                    "CompanyName": "Assistenza al Fabbricato",
                    "Config": null,
                    "DailyAllowance": "",
                    "Deductible": "",
                    "Description": "Prestazioni di assistenza e servizio 24/7 al Fabbricato quali invio di artigiani come: idraulico, elettricista, fabbro, serrandista, vetraio, sorvegliante, nei casi indicati in polizza di necessità (la compagnia eroga direttamente la prestazione non il rimborso delle spese)",
                    "Group": "ASSSISTANCE",
                    "IsBase": true,
                    "IsExtension": false,
                    "IsPremium": true,
                    "IsSellable": false,
                    "IsYour": false,
                    "IsYuor": true,
                    "LegalDefence": "",
                    "Name": "",
                    "Offer": null,
                    "Price": 0,
                    "PriceGross": 39.336,
                    "PriceNett": 35.76,
                    "SelfInsurance": "",
                    "SelfInsuranceDesc": "",
                    "Slug": "assistance",
                    "SumInsuredLimitOfIndemnity": 0,
                    "Tax": 10,
                    "Taxes": null,
                    "Type": "building",
                    "TypeOfSumInsured": "",
                    "Value": null,
                    "translation": "Assistenza"
                },
                {
                    "Assistance": "",
                    "CompanyCodec": "BI",
                    "CompanyName": "Business Interruption",
                    "Config": null,
                    "DailyAllowance": "250",
                    "Deductible": "5days",
                    "Description": "Indennizzo per il periodo di documentata inattività forzata, a seguito di un Sinistro avvenuto nel Fabbricato, che abbia danneggiato i locali e/o i Macchinar e/o le Apparecchiature Elettroniche funzionali all’attività",
                    "Group": "BUSINESS INTERRUPTTION",
                    "IsBase": false,
                    "IsExtension": false,
                    "IsPremium": true,
                    "IsSellable": false,
                    "IsYour": false,
                    "IsYuor": true,
                    "LegalDefence": "",
                    "Name": "",
                    "Offer": null,
                    "Price": 0,
                    "PriceGross": 101.959125,
                    "PriceNett": 84.09,
                    "SelfInsurance": "",
                    "SelfInsuranceDesc": "",
                    "Slug": "business-interruption",
                    "SumInsuredLimitOfIndemnity": 25000,
                    "Tax": 21.25,
                    "Taxes": null,
                    "Type": "building",
                    "TypeOfSumInsured": "firstLoss",
                    "Value": null,
                    "translation": "Business interruption"
                },
                {
                    "Assistance": "",
                    "CompanyCodec": "EL",
                    "CompanyName": "Garanzia Apparecchiature Elettroniche",
                    "Config": null,
                    "DailyAllowance": "",
                    "Deductible": "1000",
                    "Description": "Danni materiali e diretti ad apparecchiature elettroniche fisse e ad impiego mobile, causati da qualsiasi evento accidentale",
                    "Group": "ELETRONIC",
                    "IsBase": false,
                    "IsExtension": false,
                    "IsPremium": true,
                    "IsSellable": false,
                    "IsYour": false,
                    "IsYuor": false,
                    "LegalDefence": "",
                    "Name": "",
                    "Offer": null,
                    "Price": 0,
                    "PriceGross": 247.738,
                    "PriceNett": 204.32,
                    "SelfInsurance": "",
                    "SelfInsuranceDesc": "",
                    "Slug": "electronic-equipment",
                    "SumInsuredLimitOfIndemnity": 10000,
                    "Tax": 21.25,
                    "Taxes": null,
                    "Type": "building",
                    "TypeOfSumInsured": "firstLoss",
                    "Value": null,
                    "translation": "Elettronica"
                },
                {
                    "Assistance": "",
                    "CompanyCodec": "IF",
                    "CompanyName": "Incendio Fabbricato",
                    "Config": null,
                    "DailyAllowance": "",
                    "Deductible": "1000",
                    "Description": "Danni diretti a Fabbricato, causati da eventi quali: incendio, esplosione, scoppio, fulmine, conseguenti fumi, gas e vapori",
                    "Group": "FIRE",
                    "IsBase": true,
                    "IsExtension": false,
                    "IsPremium": true,
                    "IsSellable": false,
                    "IsYour": false,
                    "IsYuor": true,
                    "LegalDefence": "",
                    "Name": "",
                    "Offer": null,
                    "Price": 0,
                    "PriceGross": 84.98819999999999,
                    "PriceNett": 69.52,
                    "SelfInsurance": "",
                    "SelfInsuranceDesc": "",
                    "Slug": "building",
                    "SumInsuredLimitOfIndemnity": 120000,
                    "Tax": 22.25,
                    "Taxes": null,
                    "Type": "building",
                    "TypeOfSumInsured": "",
                    "Value": null,
                    "translation": "Fabbricato"
                },
                {
                    "Assistance": "",
                    "CompanyCodec": "IC",
                    "CompanyName": "Incendio Contenuto",
                    "Config": null,
                    "DailyAllowance": "",
                    "Deductible": "1000",
                    "Description": "Danni diretti a Contenuto (Merci, macchinari, attrezzature, arredamento), causati da eventi quali: incendio, esplosione, scoppio, fulmine, conseguenti fumi, gas e vapori",
                    "Group": "FIRE",
                    "IsBase": true,
                    "IsExtension": false,
                    "IsPremium": true,
                    "IsSellable": false,
                    "IsYour": false,
                    "IsYuor": true,
                    "LegalDefence": "",
                    "Name": "",
                    "Offer": null,
                    "Price": 0,
                    "PriceGross": 44.193374999999996,
                    "PriceNett": 36.15,
                    "SelfInsurance": "",
                    "SelfInsuranceDesc": "",
                    "Slug": "content",
                    "SumInsuredLimitOfIndemnity": 50000,
                    "Tax": 22.25,
                    "Taxes": null,
                    "Type": "building",
                    "TypeOfSumInsured": "",
                    "Value": null,
                    "translation": "Contenuto"
                },
                {
                    "Assistance": "",
                    "CompanyCodec": "EA",
                    "CompanyName": "Eventi Atmosferici",
                    "Config": null,
                    "DailyAllowance": "",
                    "Deductible": "1500",
                    "Description": "Danni da eventi atmosferici quali uragano, bufera, tempesta, grandine, vento e cose trascinate da esso, tromba d’aria, gelo, sovraccarico di neve",
                    "Group": "FIRE",
                    "IsBase": false,
                    "IsExtension": true,
                    "IsPremium": true,
                    "IsSellable": false,
                    "IsYour": false,
                    "IsYuor": true,
                    "LegalDefence": "",
                    "Name": "",
                    "Offer": null,
                    "Price": 0,
                    "PriceGross": 49.413450000000005,
                    "PriceNett": 40.42,
                    "SelfInsurance": "0.10",
                    "SelfInsuranceDesc": "10% - minimo € 1500",
                    "Slug": "atmospheric-event",
                    "SumInsuredLimitOfIndemnity": 170000,
                    "Tax": 22.25,
                    "Taxes": null,
                    "Type": "building",
                    "TypeOfSumInsured": "replacementValue",
                    "Value": null,
                    "translation": "Eventi Atmosferici"
                },
                {
                    "Assistance": "",
                    "CompanyCodec": "AC",
                    "CompanyName": "Danni d’Acqua",
                    "Config": null,
                    "DailyAllowance": "",
                    "Deductible": "1000",
                    "Description": "Danni da allagamento (eccesso o accumulo d’acqua in luogo normalmente asciutto) verificatosi all'interno del Fabbricato a seguito di formazione di ruscelli o accumulo esterno di acqua",
                    "Group": "FIRE",
                    "IsBase": true,
                    "IsExtension": true,
                    "IsPremium": true,
                    "IsSellable": false,
                    "IsYour": false,
                    "IsYuor": true,
                    "LegalDefence": "",
                    "Name": "",
                    "Offer": null,
                    "Price": 0,
                    "PriceGross": 14.828925000000002,
                    "PriceNett": 12.13,
                    "SelfInsurance": "",
                    "SelfInsuranceDesc": "",
                    "Slug": "burst-pipe",
                    "SumInsuredLimitOfIndemnity": 5000,
                    "Tax": 22.25,
                    "Taxes": null,
                    "Type": "building",
                    "TypeOfSumInsured": "firstLoss",
                    "Value": null,
                    "translation": "Danni d'acqua"
                },
                {
                    "Assistance": "",
                    "CompanyCodec": "TR",
                    "CompanyName": "Terremoto",
                    "Config": null,
                    "DailyAllowance": "",
                    "Deductible": "1000",
                    "Description": "Danni da terremoto ai beni assicurati  compresi quelli di Incendio, Esplosione, Scoppio",
                    "Group": "FIRE",
                    "IsBase": false,
                    "IsExtension": true,
                    "IsPremium": true,
                    "IsSellable": false,
                    "IsYour": false,
                    "IsYuor": false,
                    "LegalDefence": "",
                    "Name": "",
                    "Offer": null,
                    "Price": 0,
                    "PriceGross": 21.772724999999998,
                    "PriceNett": 17.81,
                    "SelfInsurance": "0.05",
                    "SelfInsuranceDesc": "5% - minimo € 1000",
                    "Slug": "earthquake",
                    "SumInsuredLimitOfIndemnity": 0.7,
                    "Tax": 22.25,
                    "Taxes": null,
                    "Type": "building",
                    "TypeOfSumInsured": "replacementValue",
                    "Value": null,
                    "translation": "Terremoto"
                },
                {
                    "Assistance": "",
                    "CompanyCodec": "RO",
                    "CompanyName": "Rottura Lastre",
                    "Config": null,
                    "DailyAllowance": "",
                    "Deductible": "1000",
                    "Description": "Spese sostenute per la sostituzione di Lastre e insegne con altre nuove eguali o equivalenti per caratteristiche, compresi i costi di trasporto ed installazione, la cui rottura sia avvenuta per cause Accidentali o imputabili a fatti di terzi.",
                    "Group": "FIRE",
                    "IsBase": false,
                    "IsExtension": true,
                    "IsPremium": true,
                    "IsSellable": false,
                    "IsYour": false,
                    "IsYuor": true,
                    "LegalDefence": "",
                    "Name": "",
                    "Offer": null,
                    "Price": 0,
                    "PriceGross": 26.418225,
                    "PriceNett": 21.61,
                    "SelfInsurance": "",
                    "SelfInsuranceDesc": "",
                    "Slug": "glass",
                    "SumInsuredLimitOfIndemnity": 2500,
                    "Tax": 22.25,
                    "Taxes": null,
                    "Type": "building",
                    "TypeOfSumInsured": "firstLoss",
                    "Value": null,
                    "translation": "Rottura Lastre"
                },
                {
                    "Assistance": "",
                    "CompanyCodec": "GU",
                    "CompanyName": "Guasti Macchine",
                    "Config": null,
                    "DailyAllowance": "",
                    "Deductible": "1500",
                    "Description": "Danni materiali e diretti, al Macchinario causati o dovuti a guasti Accidentali meccanici in genere.",
                    "Group": "FIRE",
                    "IsBase": false,
                    "IsExtension": true,
                    "IsPremium": true,
                    "IsSellable": false,
                    "IsYour": false,
                    "IsYuor": false,
                    "LegalDefence": "",
                    "Name": "",
                    "Offer": null,
                    "Price": 0,
                    "PriceGross": 71.137275,
                    "PriceNett": 58.19,
                    "SelfInsurance": "",
                    "SelfInsuranceDesc": "",
                    "Slug": "machinery-breakdown",
                    "SumInsuredLimitOfIndemnity": 5000,
                    "Tax": 22.25,
                    "Taxes": null,
                    "Type": "building",
                    "TypeOfSumInsured": "firstLoss",
                    "Value": null,
                    "translation": "Guasto Macchine"
                },
                {
                    "Assistance": "",
                    "CompanyCodec": "FE",
                    "CompanyName": "Fenomeno Elettrico",
                    "Config": null,
                    "DailyAllowance": "",
                    "Deductible": "1000",
                    "Description": "Danni originati da scariche, correnti, corto circuito ed altri fenomeni elettrici",
                    "Group": "FIRE",
                    "IsBase": false,
                    "IsExtension": true,
                    "IsPremium": true,
                    "IsSellable": false,
                    "IsYour": false,
                    "IsYuor": true,
                    "LegalDefence": "",
                    "Name": "",
                    "Offer": null,
                    "Price": 0,
                    "PriceGross": 81.40627500000001,
                    "PriceNett": 66.59,
                    "SelfInsurance": "0.10",
                    "SelfInsuranceDesc": "10% - minimo € 1000",
                    "Slug": "power-surge",
                    "SumInsuredLimitOfIndemnity": 5000,
                    "Tax": 22.25,
                    "Taxes": null,
                    "Type": "building",
                    "TypeOfSumInsured": "firstLoss",
                    "Value": null,
                    "translation": "Fenomeno Elettrico"
                },
                {
                    "Assistance": "",
                    "CompanyCodec": "ES",
                    "CompanyName": "Eventi Sociopolitici",
                    "Config": null,
                    "DailyAllowance": "",
                    "Deductible": "1000",
                    "Description": "Danni da eventi sociopolitici, quali tumulti popolari, scioperi, sommosse, atti vandalici o dolosi e atti di sabotaggio",
                    "Group": "FIRE",
                    "IsBase": false,
                    "IsExtension": true,
                    "IsPremium": true,
                    "IsSellable": false,
                    "IsYour": false,
                    "IsYuor": true,
                    "LegalDefence": "",
                    "Name": "",
                    "Offer": null,
                    "Price": 0,
                    "PriceGross": 6.124725,
                    "PriceNett": 5.01,
                    "SelfInsurance": "0.05",
                    "SelfInsuranceDesc": "5% - minimo € 1000",
                    "Slug": "sociopolitical-event",
                    "SumInsuredLimitOfIndemnity": 136000,
                    "Tax": 22.25,
                    "Taxes": null,
                    "Type": "building",
                    "TypeOfSumInsured": "",
                    "Value": null,
                    "translation": "Eventi Sociopolitici"
                },
                {
                    "Assistance": "",
                    "CompanyCodec": "AT",
                    "CompanyName": "Atti di Terrorismo",
                    "Config": null,
                    "DailyAllowance": "",
                    "Deductible": "1000",
                    "Description": "Danni da atti di terrorismo",
                    "Group": "FIRE",
                    "IsBase": false,
                    "IsExtension": true,
                    "IsPremium": true,
                    "IsSellable": false,
                    "IsYour": false,
                    "IsYuor": true,
                    "LegalDefence": "",
                    "Name": "",
                    "Offer": null,
                    "Price": 0,
                    "PriceGross": 0,
                    "PriceNett": 0,
                    "SelfInsurance": "0.05",
                    "SelfInsuranceDesc": "5% - minimo € 1000",
                    "Slug": "terrorism",
                    "SumInsuredLimitOfIndemnity": 0.5,
                    "Tax": 22.25,
                    "Taxes": null,
                    "Type": "building",
                    "TypeOfSumInsured": "",
                    "Value": null,
                    "translation": "Atti di Terrorismo"
                },
                {
                    "Assistance": "",
                    "CompanyCodec": "AG",
                    "CompanyName": "Allagamento",
                    "Config": null,
                    "DailyAllowance": "",
                    "Deductible": "1000",
                    "Description": "Danni da fuoriuscita di acqua condotta (ingorghi, trabocchi, rotture accedentali), e, se assicurato il fabbricato, ricerca ripristino e riparazione del danno",
                    "Group": "FIRE",
                    "IsBase": false,
                    "IsExtension": true,
                    "IsPremium": true,
                    "IsSellable": false,
                    "IsYour": false,
                    "IsYuor": true,
                    "LegalDefence": "",
                    "Name": "",
                    "Offer": null,
                    "Price": 0,
                    "PriceGross": 15.171225,
                    "PriceNett": 12.41,
                    "SelfInsurance": "0.05",
                    "SelfInsuranceDesc": "5% - minimo € 1000",
                    "Slug": "water-damage",
                    "SumInsuredLimitOfIndemnity": 0.7,
                    "Tax": 22.25,
                    "Taxes": null,
                    "Type": "building",
                    "TypeOfSumInsured": "replacementValue",
                    "Value": null,
                    "translation": "Allagamento"
                },
                {
                    "Assistance": "",
                    "CompanyCodec": "RF",
                    "CompanyName": "RC Fabbricato",
                    "Config": null,
                    "DailyAllowance": "",
                    "Deductible": "1000",
                    "Description": "Danni a persone o cose, verificatosi in relazione alla proprietà e conduzione del Fabbricato e delle eventuali parti comuni a esso riferite/collegate",
                    "Group": "RCF",
                    "IsBase": true,
                    "IsExtension": false,
                    "IsPremium": true,
                    "IsSellable": false,
                    "IsYour": false,
                    "IsYuor": true,
                    "LegalDefence": "",
                    "Name": "",
                    "Offer": null,
                    "Price": 0,
                    "PriceGross": 123.59474999999999,
                    "PriceNett": 101.1,
                    "SelfInsurance": "",
                    "SelfInsuranceDesc": "",
                    "Slug": "property-owners-liability",
                    "SumInsuredLimitOfIndemnity": 1000000,
                    "Tax": 22.25,
                    "Taxes": null,
                    "Type": "building",
                    "TypeOfSumInsured": "",
                    "Value": null,
                    "translation": "Responsabilità Civile Fabbricato"
                },
                {
                    "Assistance": "",
                    "CompanyCodec": "RT",
                    "CompanyName": "Ricorso Terzi da Incendio",
                    "Config": null,
                    "DailyAllowance": "",
                    "Deductible": "1000",
                    "Description": "Responsabilità per danni materiali e diretti arrecati alle cose di terzi in seguito a Incendio, Esplosione o Scoppio del Fabbricato e/o Contenuto, qualora assicurati, anche quando il Fabbricato lo è nella forma di Rischio Locativo",
                    "Group": "RT",
                    "IsBase": true,
                    "IsExtension": false,
                    "IsPremium": true,
                    "IsSellable": false,
                    "IsYour": false,
                    "IsYuor": true,
                    "LegalDefence": "",
                    "Name": "",
                    "Offer": null,
                    "Price": 0,
                    "PriceGross": 111.7854,
                    "PriceNett": 91.44,
                    "SelfInsurance": "0.10",
                    "SelfInsuranceDesc": "10% - minimo € 1000",
                    "Slug": "third-party-recourse",
                    "SumInsuredLimitOfIndemnity": 150000,
                    "Tax": 22.25,
                    "Taxes": null,
                    "Type": "building",
                    "TypeOfSumInsured": "firstLoss",
                    "Value": null,
                    "translation": "Ricorso Terzi incendio"
                },
                {
                    "Assistance": "",
                    "CompanyCodec": "FU",
                    "CompanyName": "Garanzia Furto, Rapina ed Estorsione",
                    "Config": null,
                    "DailyAllowance": "",
                    "Deductible": "1000",
                    "Description": "Danni subiti da furto, rapina o estorsione, inclusi i guasti e gli atti vandalici commessi dai ladri",
                    "Group": "THEFT",
                    "IsBase": false,
                    "IsExtension": false,
                    "IsPremium": true,
                    "IsSellable": false,
                    "IsYour": false,
                    "IsYuor": false,
                    "LegalDefence": "",
                    "Name": "",
                    "Offer": null,
                    "Price": 0,
                    "PriceGross": 66.87075,
                    "PriceNett": 54.7,
                    "SelfInsurance": "",
                    "SelfInsuranceDesc": "10% - minimo € 1000",
                    "Slug": "theft",
                    "SumInsuredLimitOfIndemnity": 5000,
                    "Tax": 22.25,
                    "Taxes": null,
                    "Type": "building",
                    "TypeOfSumInsured": "firstLoss",
                    "Value": null,
                    "translation": "Furto, rapina, estorsione"
                }
            ]
        }
    ],
    "statements": [
        {
            "title": "1) Scelta firma elettronica",
            "answer": true,
            "questions": [
                {
                    "question": "Si prende e si dà atto tra le Parti che, il contratto viene sottoscritto con Firma Elettronica Avanzata, redatto in un unico esemplare. Pertanto:",
                    "isbold": true,
                    "indent": false
                },
                {
                    "question": "dichiaro di aver ricevuto, preso visione, conoscere ed accettare le \"Condizioni Generali di Servizio per l'utilizzazione della Firma Elettronica Avanzata” prevista da Wopta e l’annessa “Scheda Tecnica Illustrativa”",
                    "isbold": false,
                    "indent": false
                },
                {
                    "question": "confermo la veridicità dei dati forniti, la titolarità del numero di cellulare e dell’indirizzo mail, acconsentendo al trattamento di tali dati per questa specifica finalità;",
                    "isbold": false,
                    "indent": false
                },
                {
                    "question": "dichiaro altresì di avere titolo a richiedere l'attivazione e l’uso del Servizio per la sottoscrizione del presente contratto",
                    "isbold": false,
                    "indent": false
                }
            ]
        },
        {
            "title": "2) Dichiarazioni sul rischio ",
            "answer": true,
            "questions": [
                {
                    "question": "Premesso di essere a conoscenza che le dichiarazioni non veritiere, inesatte o reticenti, da me rese, possono compromettere il diritto alla prestazione (come da art. 1892, 1893, 1894 c.c.), ai fini dell’efficacia delle garanzie ",
                    "isbold": false,
                    "indent": false
                },
                {
                    "question": "dichiaro che:",
                    "isbold": false,
                    "indent": false
                },
                {
                    "question": "l’azienda assicurata e/o gli immobili assicurati, rispondono ai requisiti indicati all’art. 9 – “requisiti di assicurabilità” delle condizioni di assicurazione;",
                    "isbold": false,
                    "indent": false
                },
                {
                    "question": "che, sui medesimi rischi assicurati con la presente Polizza, nel triennio precedente:",
                    "isbold": false,
                    "indent": false
                },
                {
                    "question": "NON vi sono state coperture assicurative annullate dall’assicuratore;",
                    "isbold": false,
                    "indent": false
                },
                {
                    "question": "NON si sono verificati eventi dannosi di importo liquidato superiore a 1.000 €",
                    "isbold": false,
                    "indent": false
                },
                {
                    "question": "al momento della stipula di questa Polizza NON ha ricevuto comunicazioni, richieste e notifiche che possano configurare un sinistro relativo alle garanzie assicurate e di non essere a conoscenza di eventi o circostanze che possano dare origine ad una richiesta di risarcimento.",
                    "isbold": false,
                    "indent": false
                }
            ]
        },
        {
            "title": "3) SCELTA COMUNICAZIONI VIA MAIL E ACCETTAZIONE POLIZZA",
            "answer": true,
            "questions": [
                {
                    "question": "Ho scelto la ricezione della seguente documentazione via e-mail al seguente indirizzo: luca.barbieri@wopta.it. Sono a conoscenza che, anche le future comunicazioni avverranno con questo mezzo e che qualora volessi modificare questa mia scelta potrò farlo scrivendo a Global Assistance, con le modalità previste nelle Condizioni di Assicurazione.",
                    "isbold": false,
                    "indent": false
                },
                {
                    "question": "Confermo quindi di aver ricevuto e preso visione, prima della conclusione del contratto:",
                    "isbold": false,
                    "indent": false
                },
                {
                    "question": "degli Allegati 3, 4 e 4-ter, di cui al Regolamento IVASS n. 40/2018, relativi agli obblighi informativi e di comportamento dell’Intermediario, inclusa l’informativa privacy dell’intermediario (ai sensi dell’art. 13 del regolamento UE n. 2016/679);",
                    "isbold": false,
                    "indent": false
                },
                {
                    "question": "del Set informativo, identificato dal modello GA02.0922, contenente: 1) documento informativo per i prodotti assicurativi danni (DIP Danni) e documento informativo precontrattuale aggiuntivo per i prodotti assicurativi danni (DIP Aggiuntivo danni) cui al Regolamento IVASS n. 41/2018; 2) Condizioni di Assicurazione comprensive di Glossario, che dichiaro altresì di conoscere ed accettare.",
                    "isbold": false,
                    "indent": false
                }
            ]
        },
        {
            "title": "4) CLAUSOLE DA APPROVARE IN MODO SPECIFICO",
            "answer": true,
            "questions": [
                {
                    "question": "Ai sensi degli artt. 1341 e 1342 Codice Civile, dichiaro di approvare in modo specifico, le disposizioni indicate nelleCondizioni di Assicurazione con particolare riguardo agli articoli dei seguenti capitoli:",
                    "isbold": false,
                    "indent": false
                },
                {
                    "question": "Art. 1.3 Revisione del Premio e/o condizioni; Art. 4 Foro competente – Arbitrato;- Art. 15.1 Recesso in caso di sinistro; Art. 64 Garanzia “Difesa Penale” – Cosa è assicurato; Art. 70 Garanzia “Difesa Penale, Civile e Circolazione” – Cosa è assicurato; Art. 73.3 Estorsione Cyber; Art. 73.9 Cyber-Crime; Art. 95.3 Procedura per la valutazione del danno; Art. 95.4 Mandato dei Periti; Art. 95.5 Determinazione del danno e valore dei beni assicurati; Art. 95.13 Titolarità dei diritti nascenti dalla Polizza; Art. 97.4 Procedura per la valutazione del danno; Art. 97.6 Operazioni peritali; Art. 97.7 determinazione del danno e valore dei beni assicurati; Art. 97.14 Titolarità dei diritti nascenti dalla Polizza; Art. 98.1 Cosa fare al momento del sinistro; Art. 100 Sezione F – Cyber-Risk; Art. 100.2 Gestione delle richieste di Risarcimento; Art. 100.4 Titolarità dei diritti nascenti dalla Polizza; Art. 102.4 Procedura per la valutazione del danno; Art. 102.6 Operazioni peritali; Art. 102.7 Determinazione del danno; Art. 102.10 Titolarità dei diritti nascenti dalla Polizza.",
                    "isbold": false,
                    "indent": false
                }
            ]
        },
        {
            "title": "5) PER NOI QUESTA POLIZZA FA AL CASO TUO",
            "answer": true,
            "questions": [
                {
                    "question": "Hai effettuato dichiarazioni relative al rischio da assicurare e scelto prestazioni e garanzie tra quelle proposte. Sulla base di tali dichiarazioni, delle tue esigenze e richieste, le soluzioni assicurative individuate e assolte dalle coperture assicurative risultano le seguenti:",
                    "isbold": false,
                    "indent": false
                },
                {
                    "question": "1. Assicurare l’Attività aziendale svolta nell’ubicazione indicata in Polizza, per rischi relativi a:",
                    "isbold": true,
                    "indent": false
                },
                {
                    "question": "1.1. Danni involontariamente causati a terzi per danni a cose o persone di cui sia responsabile a termini di legge (R.C.T.). La garanzia include, ma non si limita a questi, i danni: a veicoli di terzi e prestatori di lavoro; a cose in consegna e custodia; a cose nell'ambito di esecuzione dei lavori; a cose di terzi sollevate, caricate, scaricate, movimentate, trasportate o rimorchiate; a mezzi di trasporto sotto carico e scarico; da interruzione o sospensione di attività di terzi; da smercio; da committenza autoveicoli; da responsabilità civile personale addetti; da attività di commercio ambulante; da lavori presso terzi.",
                    "isbold": false,
                    "indent": true
                },
                {
                    "question": "1.2. Responsabilità per: la rivalsa INAIL per gli infortuni sul lavoro subiti dai prestatori di lavoro; morte; e lesioni personali dalle quali sia derivata un'invalidità permanente ai sensi del codice civile, incluse le malattie professionali",
                    "isbold": false,
                    "indent": true
                },
                {
                    "question": "1.3. Difesa penale: per reati di natura colposa o contravvenzionale, inclusi i casi di sicurezza aziendale da D. Lgs. 81/08 e D. Lgs. 106/09, D. Lgs. 193/07, D. Lgs. 152/06, D. Lgs. 101/18, D. Lgs. 231/01",
                    "isbold": false,
                    "indent": true
                },
                {
                    "question": "2. Assicurare gli asset relativi alla Sede dell’Attività indicata in Polizza, in Galleria del Corso, 1 per rischi relativi a: ",
                    "isbold": true,
                    "indent": false
                },
                {
                    "question": "2.1. Danni diretti a Fabbricato e Contenuto, causati da eventi quali: incendio, esplosione, scoppio, fulmine, conseguenti fumi, gas e vapori; terremoto; allagamento; guasto macchine; rottura di lastre; danni da fuoriuscita di acqua condotta, e, se assicurato il fabbricato, ricerca ripristino e riparazione del danno; eventi atmosferici quali uragano, bufera, tempesta, grandine, vento e cose trascinate da esso, tromba d’aria, gelo, sovraccarico di neve; eventi sociopolitici, quali tumulti popolari, scioperi, sommosse, atti vandalici o dolosi e atti di sabotaggio; scariche, correnti, corto circuito ed altri fenomeni elettrici; atti di terrorismo;",
                    "isbold": false,
                    "indent": true
                },
                {
                    "question": "2.2. Responsabilità per danni materiali e diretti arrecati alle cose di terzi in seguito a Incendio, Esplosione o Scoppio del Fabbricato e/o Contenuto, qualora assicurati, anche quando il Fabbricato lo è nella forma di Rischio Locativo;",
                    "isbold": false,
                    "indent": true
                },
                {
                    "question": "2.3. Danni a persone o cose, verificatosi in relazione alla proprietà e conduzione del Fabbricato e delle eventuali parti comuni a esso riferite/collegate;",
                    "isbold": false,
                    "indent": true
                },
                {
                    "question": "2.4. Danni subiti da furto, rapina o estorsione, inclusi i guasti e gli atti vandalici commessi dai ladri;",
                    "isbold": false,
                    "indent": true
                },
                {
                    "question": "2.5. Perdite pecuniarie per il periodo di documentata inattività forzata, a seguito di un sinistro avvenuto nel Fabbricato, che abbia danneggiato i locali e/o i macchinari e/o le apparecchiature elettroniche funzionali all’Attività, indennizzabile ai sensi degli eventi garantiti alla lettera a) che precede;",
                    "isbold": false,
                    "indent": true
                },
                {
                    "question": "2.6. Danni materiali e diretti ad apparecchiature elettroniche fisse e ad impiego mobile, causati da qualsiasi evento accidentale, qualunque ne sia la causa, non espressamente escluso;",
                    "isbold": false,
                    "indent": true
                },
                {
                    "question": "2.7. Prestazioni di assistenza e servizio 24/7 al Fabbricato quali invio di artigiani come: idraulico, elettricista, fabbro, serrandista, vetraio, sorvegliante, nei casi indicati in polizza di necessità (la compagnia eroga direttamente la prestazione non il rimborso delle spese);",
                    "isbold": false,
                    "indent": true
                },
                {
                    "question": "La Polizza prevede, in relazione a tutte le garanzie che precedono, l’applicazione di Scoperti, Franchigie, Limiti di indennizzo ed esclusioni, meglio riportate nelle Condizioni Generali di Assicurazione.",
                    "isbold": true,
                    "indent": false
                },
                {
                    "question": "Al livello di Franchigia scelto, MEDIO-BASSO (1000€) corrisponde per ogni garanzia, nella Tabella “Scoperti e Franchigie”, il dettaglio di tutti gli Scoperti e Franchigie applicabili in caso di Sinistro, di cui l’importo qui indicato, costituisce il minimo se non diversamente specificato. Tali importi sono stati da te valutati in linea con la capacità finanziaria di sostenere in proprio tale livello di danno e rischio.",
                    "isbold": false,
                    "indent": false
                },
                {
                    "question": "Con la seguente sottoscrizione dichiari che quanto precede corrisponde alle informazioni ottenute dall’intermediario, sia attraverso i documenti resi disponibili e/o inviati che nelle pagine web del sito wopta.it.",
                    "isbold": false,
                    "indent": false
                }
            ]
        }
    ]
}`

}
