package query_builder

type QueryBuilder interface {
	BuildQuery(map[string]string) (string, map[string]interface{})
}
