package companydata

type fn func(interface{})interface{}
func GetMapFx(name string,value interface{})interface{}{

res:= map[string]fn{
    "test": Test,
   
}
	return res[name](value)
}
func Test(s interface{})interface{} {
	return ""
}