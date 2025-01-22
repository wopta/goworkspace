package pkg

type QueryBuilder interface {
	Build(map[string]string) (string, map[string]interface{})
}
