package lib

import "reflect"

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

			if typeValue.Type.Kind() == reflect.Struct {
				FieldNames(fieldValue.Interface(), m)

			}
		}
	}
}
