package renew

import "github.com/wopta/goworkspace/models"

type RenewResp struct {
	Success []RenewReport `json:"success"`
	Failure []RenewReport `json:"failure"`
}

type RenewReport struct {
	Policy       models.Policy        `json:"policy"`
	Transactions []models.Transaction `json:"transactions"`
	Error        error                `json:"error,omitempty"`
}
