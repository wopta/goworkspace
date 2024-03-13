package models

import "cloud.google.com/go/bigquery"

type NetworkTreeElement struct {
	RootUid       string                `json:"rootUid" bigquery:"rootUid"`
	ParentUid     string                `json:"parentUid" bigquery:"parentUid"`
	NodeUid       string                `json:"nodeUid" bigquery:"nodeUid"`
	Name          string                `json:"name" bigquery:"name"`
	AbsoluteLevel int                   `json:"absoluteLevel" bigquery:"absoluteLevel"`
	RelativeLevel int                   `json:"relativeLevel" bigquery:"relativeLevel"`
	CreationDate  bigquery.NullDateTime `json:"-" bigquery:"creationDate"`
	Ancestors     []NetworkTreeElement  `json:"ancestors,omitempty" bigquery:"-"`
	Children      []NetworkTreeElement  `json:"children,omitempty" bigquery:"-"`
}
