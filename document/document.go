package document

import (
	"net/http"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/props"
	"gitlab.dev.wopta.it/goworkspace/lib"
)

type DocumentResp struct {
	LinkGcs  string `json:"linkGcs"`
	Filename string `json:"filename"`
}

var documentRoutes []lib.Route = []lib.Route{
	{
		Route:   "/generate/mup/{nodeUid}",
		Handler: lib.ResponseLoggerWrapper(GenerateMupFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAll},
	},
}

func init() {
	log.Println("INIT Document")
	functions.HTTP("Document", Document)
}

func Document(w http.ResponseWriter, r *http.Request) {
	router := lib.GetRouter("document", documentRoutes)
	router.ServeHTTP(w, r)
}

var signatureID int

type Kv struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type DocumentResponse struct {
	EnvelopSignId string `json:"envelopSignId"`
	LinkGcs       string `json:"linkGcs"`
	Bytes         string `json:"bytes"`
	FullPath      string
}

type DodumentData struct {
	Class        string `json:"class"`
	CoverageType string `json:"coverageType"`
	FiscalCode   string `json:"fiscalCode"`
	Name         string `json:"name"`
	Price        int64  `json:"price"`
	PriceNett    int64  `json:"priceNett"`
	Surname      string `json:"surname"`
	Work         string `json:"work"`
	WorkType     string `json:"workType"`
	Coverages    []struct {
		Deductible                 string `json:"deductible"`
		Name                       string `json:"name"`
		Price                      int64  `json:"price"`
		PriceNett                  int64  `json:"priceNett"`
		SelfInsurance              string `json:"selfInsurance"`
		SumInsuredLimitOfIndemnity int64  `json:"sumInsuredLimitOfIndemnity"`
	} `json:"coverages"`
}

type Skin struct {
	PrimaryColor         color.Color
	SecondaryColor       color.Color
	LineColor            color.Color
	TextColor            color.Color
	TitleColor           color.Color
	RowHeight            float64
	rowtableHeight       float64
	LineHeight           float64
	Size                 float64
	SizeTitle            float64
	RowTitleHeight       float64
	TableHeight          float64
	rowtableHeightMin    float64
	DynamicHeightMin     int
	CharForRow           int
	DynamicHeightDiv     float64
	MagentaTextLeft      props.Text
	WhiteTextCenter      props.Text
	MagentaBoldtextRight props.Text
	MagentaBoldtextLeft  props.Text
	MagentatextRight     props.Text
	MagentatextLeft      props.Text
	NormaltextLeft       props.Text
	NormaltextLeftBlack  props.Text
	BoldtextLeft         props.Text
	NormaltextRight      props.Text
	NormaltextLeftExt    props.Text
	TitletextLeft        props.Text
	TitletextRight       props.Text
	TitletextCenter      props.Text
	TitleBoldtextRight   props.Text
	TitleBoldtextLeft    props.Text
	Line                 props.Line
}
