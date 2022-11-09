package models

type Policy struct {
	ID            string `firestore:"id,omitempty" json:"id,omitempty"`
	IdSign        string `firestore:"idSign,omitempty" json:"idSign,omitempty"`
	Uid           string `firestore:"uid,omitempty" json:"uid,omitempty"`
	Number        string `firestore:"number,omitempty" json:"number,omitempty"`
	NumberCompany string `firestore:"numberCompany,omitempty" json:"numberCompany,omitempty"`
	Status        string `firestore:"status ,omitempty" json:"status ,omitempty"`
	Company       string `firestore:"company,omitempty" json:"company,omitempty"`
	Name          string `firestore:"name,omitempty" json:"name,omitempty"`
	StartDate     string `firestore:"startDate,omitempty" json:"startDate,omitempty"`
	EndDate       string `firestore:"endDate,omitempty" json:"endDate,omitempty"`
	Payment       string `firestore:"payment,omitempty" json:"payment,omitempty"`
	PaymentType   string `firestore:"paymentType,omitempty" json:"paymentType,omitempty"`
	PaymentSplit  string `firestore:"paymentSplit,omitempty" json:"paymentSplit,omitempty"`
	IsPay         bool   `firestore:"isPay,omitempty" json:"isPay,omitempty"`
	IsSign        bool   `firestore:"isSign,omitempty" json:"isSign,omitempty"`
	CoverageType  string `firestore:"coverageType,omitempty" json:"coverageType,omitempty"`
	Voucher       string `firestore:"voucher,omitempty" json:"voucher,omitempty"`
	Channel       string `firestore:"channel,omitempty" json:"channel,omitempty"`
	Covenant      string `firestore:"covenant,omitempty" json:"covenant,omitempty"`
	TaxAmount     int64  `firestore:"taxAmount,omitempty" json:"taxAmount,omitempty"`
	PriceNett     int64  `firestore:"priceNett,omitempty" json:"priceNett,omitempty"`
	PriceGross    int64  `firestore:"priceGross,omitempty" json:"priceGross   ,omitempty"`
	Contractor    *User  `firestore:"contractor,omitempty" json:"contractor,omitempty"`
	DocumentName  string `firestore:"documentName,omitempty" json:"documentName,omitempty"`

	Assets []Asset `firestore:"guarantees,omitempty" json:"guarantees,omitempty"`
}
