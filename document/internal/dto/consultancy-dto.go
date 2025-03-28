package dto

type consultancyDTO struct {
	Price      numeric
}

func newConsultacyDTO() *consultancyDTO {
	return &consultancyDTO{
		Price: newNumeric(),
	}
}
