package enrich

import (
	"net/http"
	"strings"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"github.com/go-chi/chi/v5"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"gitlab.dev.wopta.it/goworkspace/lib"
)

func AtecoFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		ricAteco   []byte
		enrichByte []byte
	)

	log.AddPrefix("AtecoFx")
	defer log.PopPrefix()

	log.Println("Handler start -----------------------------------------------")

	ateco := chi.URLParam(r, "ateco")
	log.Println(ateco)

	ricAteco = lib.GetFilesByEnv("data/rules/Riclassificazione_Ateco.csv")

	df := lib.CsvToDataframe(ricAteco)
	fil := df.Filter(
		dataframe.F{Colidx: 5, Colname: "Codice Ateco 2007", Comparator: series.Eq, Comparando: ateco},
	)
	log.Println("filtered row", fil.Nrow())
	log.Println("filtered col", fil.Ncol())

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

	log.Println("Handler end -------------------------------------------------")

	return "{\"ateco\":" + string(enrichByte) + "}", nil, nil
}
