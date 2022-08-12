package models

type ProfileAllriskJson struct {
	Vat              int64  `json:"vat"`
	SquareMeters     int64  `json:"squareMeters"`
	IsBuildingOwner  bool   `json:"isBuildingOwner"`
	Revenue          int64  `json:"revenue"`
	Address          string `json:"Address"`
	Ateco            string `json:"Ateco"`
	BusinessSector   string `json:"businessSector"`
	BuildingType     string `json:"buildingType"`
	BuildingMaterial string `json:"buildingMaterial"`
	BuildingYear     string `json:"buildingYear"`
	Employer         int64  `json:"employer"`
	IsAllarm         bool   `json:"isAllarm"`
	Floor            int64  `json:"floor"`
	IsPRA            bool   `json:"isPRA"`
	Costruction      string `json:"costruction"`
	IsHolder         bool   `json:"isHolder"`
	WhatToSay        string `json:"whatToSay"`
}
