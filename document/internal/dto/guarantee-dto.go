package dto

type GuaranteeDTO struct {
	Description                string
	SumInsuredLimitOfIndemnity float64
	LimitOfIndemnity           float64
	SumInsured                 float64
}

func newGuaranteeDTO() *GuaranteeDTO {
	return &GuaranteeDTO{
		Description:                emptyField,
		SumInsuredLimitOfIndemnity: 0,
		LimitOfIndemnity:           0,
		SumInsured:                 0,
	}
}
