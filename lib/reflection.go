package lib

import (
	"encoding/json"
	"fmt"

	"reflect"
	"strings"

	dynamicstruct "github.com/ompluscator/dynamic-struct"
	models "github.com/wopta/goworkspace/models"
)

func JsonToStruct(req []byte) {
	//lib.Files("")
	var profileAllriskJson models.ProfileAllriskJson
	var result map[string]interface{}
	//req, err := io.ReadAll(r.Body)
	json.Unmarshal(req, &result)
	ds := dynamicstruct.NewStruct()

	fmt.Println(result)
	fmt.Println(profileAllriskJson)
	keys := reflect.ValueOf(result).MapKeys()
	for i := 0; i < len(result); i++ {

		fmt.Println(strings.Title(keys[i].String()))
		fmt.Println(result[keys[i].String()])
		fmt.Println("ds.AddField")
		var d = `default:"` + fmt.Sprintf("%v", result[keys[i].String()]) + `"`
		fmt.Println(d)
		ds.AddField(strings.Title(keys[i].String()), result[keys[i].String()], "")
	}

	//ds.AddField()
	fmt.Println("ds.Build().New()")
	instance := ds.Build().New()
	fmt.Println(instance)
	ps := reflect.ValueOf(instance)
	for i := 0; i < len(result); i++ {
		fmt.Println(ps.Elem())
		fmt.Println(ps.Elem().FieldByName(strings.Title(keys[i].String())))
		ps.Elem().FieldByName(strings.Title(keys[i].String())).Set(reflect.ValueOf(result[keys[i].String()]))

	}

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	//err = json.NewDecoder(r.Body).Decode(&instance)

}
