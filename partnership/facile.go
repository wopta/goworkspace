package partnership

import (
	"github.com/golang-jwt/jwt/v4"
)

type FacileClaims struct {
	CustomerName       string `json:"customerName"`
	CustomerFamilyName string `json:"customerFamilyName"`
	CustomerBirthDate  string `json:"customerBirthDate"`
	Gender             string `json:"gender"`
	Email              string `json:"email"`
	Mobile             string `json:"mobile"`
	IsSmoker           bool   `json:"isSmoker"`
	InsuredCapital     int    `json:"insuredCapital"`
	Duration           int    `json:"duration"`
	jwt.RegisteredClaims
}
