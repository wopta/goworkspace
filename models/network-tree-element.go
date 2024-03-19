package models

import "cloud.google.com/go/bigquery"

type NetworkTreeElement struct {
	RootUid       string                `json:"rootUid" bigquery:"rootUid"`
	ParentUid     string                `json:"parentUid" bigquery:"parentUid"`
	NodeUid       string                `json:"nodeUid" bigquery:"nodeUid"`
	Name          string                `json:"name" bigquery:"name"`
	AbsoluteLevel int                   `json:"-" bigquery:"-"`
	RelativeLevel int                   `json:"relativeLevel" bigquery:"relativeLevel"`
	CreationDate  bigquery.NullDateTime `json:"-" bigquery:"creationDate"`
	Ancestors     []NetworkTreeElement  `json:"ancestors,omitempty" bigquery:"-"`
	Children      []NetworkTreeElement  `json:"children,omitempty" bigquery:"-"`
}

func (nte *NetworkTreeElement) ToNetworkTreeRelation() NetworkTreeRelation {
	return NetworkTreeRelation{
		RootUid:       nte.RootUid,
		ParentUid:     nte.ParentUid,
		NodeUid:       nte.NodeUid,
		RelativeLevel: nte.RelativeLevel,
		CreationDate:  nte.CreationDate,
	}
}

type NetworkTreeRelation struct {
	RootUid       string                `json:"rootUid" bigquery:"rootUid"`
	ParentUid     string                `json:"parentUid" bigquery:"parentUid"`
	NodeUid       string                `json:"nodeUid" bigquery:"nodeUid"`
	RelativeLevel int                   `json:"relativeLevel" bigquery:"relativeLevel"`
	CreationDate  bigquery.NullDateTime `json:"-" bigquery:"creationDate"`
}

func (ntr *NetworkTreeRelation) ToNetworkTreeElement() NetworkTreeElement {
	return NetworkTreeElement{
		RootUid:       ntr.RootUid,
		ParentUid:     ntr.ParentUid,
		NodeUid:       ntr.NodeUid,
		RelativeLevel: ntr.RelativeLevel,
		CreationDate:  ntr.CreationDate,
	}
}
