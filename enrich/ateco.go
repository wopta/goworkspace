package enrich

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	lib "github.com/wopta/goworkspace/lib"
)

func Ateco(w http.ResponseWriter, ateco string) {
	// Set CORS headers for the main request.
	var (
		ricAteco []byte
	)

	w.Header().Set("Access-Control-Allow-Methods", "GET")
	log.Println("Ateco")
	w.Header().Set("Content-Type", "application/json")
	switch os.Getenv("env") {
	case "local":

		ricAteco = lib.ErrorByte(ioutil.ReadFile("function-data/data/rules/Riclassificazione_Ateco.csv"))
	case "dev":

		ricAteco = lib.GetFromStorage("function-data", "data/rules/Riclassificazione_Ateco.csv", "")
	case "prod":

		ricAteco = lib.GetFromStorage("core-350507-function-data", "data/rules/Riclassificazione_Ateco.csv", "")
	default:

	}
	df := lib.CsvToDataframe(ricAteco)
	fil := df.Filter(
		dataframe.F{Colidx: 5, Colname: "Codice Ateco 2007", Comparator: series.Eq, Comparando: ateco},
	)
	log.Println("filtered row", fil.Nrow())
	log.Println("filtered col", fil.Ncol())
	var enrichByte []byte

	if fil.Nrow() > 0 {
		enrichByte = []byte(`{	"atecoMacro":"` + strings.ToUpper(fil.Elem(0, 0).String()) + `",
		"atecoSub":"` + strings.ToUpper(fil.Elem(0, 1).String()) + `",
		"atecoDesc":"` + strings.ToUpper(fil.Elem(0, 2).String()) + `",
		"businessSector":"` + strings.ToUpper(fil.Elem(0, 3).String()) + `",
		"fire":"` + strings.ToUpper(fil.Elem(0, 14).String()) + `",
		"fireLow500k":"` + strings.ToUpper(fil.Elem(0, 5).String()) + `",
		"fireUp500k":"` + strings.ToUpper(fil.Elem(0, 6).String()) + `",
		"theft":"` + strings.ToUpper(fil.Elem(0, 15).String()) + `",
		"thefteLow500k ":"` + strings.ToUpper(fil.Elem(0, 8).String()) + `",
		"theftUp500k":"` + strings.ToUpper(fil.Elem(0, 9).String()) + `",
		"rct":"` + strings.ToUpper(fil.Elem(0, 16).String()) + `",
		"rco":"` + strings.ToUpper(fil.Elem(0, 17).String()) + `",
		"rcoProd":"` + strings.ToUpper(fil.Elem(0, 18).String()) + `",
		"rcVehicle":"` + strings.ToUpper(fil.Elem(0, 19).String()) + `",
		"rcpo":"` + strings.ToUpper(fil.Elem(0, 20).String()) + `",
		"rcp12":"` + strings.ToUpper(strings.ToUpper(fil.Elem(0, 21).String())) + `",
		"rcp2008":"` + strings.ToUpper(fil.Elem(0, 22).String()) + `",
		"damageTheft":"` + strings.ToUpper(fil.Elem(0, 23).String()) + `",
		"damageThing":"` + strings.ToUpper(fil.Elem(0, 24).String()) + `",
		"rcCostruction":"` + strings.ToUpper(fil.Elem(0, 25).String()) + `",
		"eletronic":"` + strings.ToUpper(fil.Elem(0, 27).String()) + `",
		"machineFaliure":"` + strings.ToUpper(fil.Elem(0, 28).String()) + `"}`)
	} else {
		enrichByte = []byte(`{}`)
	}

	log.Println(string(enrichByte))

	fmt.Fprintf(w, "{\"ateco\":"+string(enrichByte)+"}")
}
