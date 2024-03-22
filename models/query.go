package models

type Query struct {
	Field  string        `json:"field"`
	Op     string        `json:"op"`
	Value  interface{}   `json:"value"`
	Type   string        `json:"type"`
	Values []interface{} `json:"values"`
}
