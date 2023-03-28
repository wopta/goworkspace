package models

type WiseLocation struct {
	Note             string      `json:"txNoteUbicazione"`
	ConstructionYear int         `json:"nAnnoCostruzione"`
	BuildingTypeCode string      `json:"cdTipoFabbricato"`
	BuildingType     string      `json:"txTipoFabbricato"`
	LayoutTypeCode   string      `json:"cdTipoRipartizione"`
	LayoutType       string      `json:"txTipoRipartizione"`
	ActivityCode     string      `json:"cdAttivita"`
	Activity         string      `json:"txAttivita"`
	IsMainAddress    bool        `json:"bAbitazionePrincipale"`
	PropertyValue    float64     `json:"nValoreImmobile"`
	Address          WiseAddress `json:"indirizzo"`
}