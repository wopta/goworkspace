module github.com/wopta/goworkspace/cmd

go 1.16

replace (
	github.com/wopta/goworkspace/cmd => ./
)

require (
	github.com/GoogleCloudPlatform/functions-framework-go v1.6.1
	github.com/wopta/goworkspace/broker v0.0.0-20230104102617-e971cb02bc2a
	github.com/wopta/goworkspace/document v0.0.0-20230104102617-e971cb02bc2a
	github.com/wopta/goworkspace/enrich v0.0.0-20230104102617-e971cb02bc2a
	github.com/wopta/goworkspace/mail v0.0.0-20230104102617-e971cb02bc2a
	github.com/wopta/goworkspace/quote v0.0.0-20230104102617-e971cb02bc2a
	github.com/wopta/goworkspace/rules v0.0.0-20230104102617-e971cb02bc2a
	//github.com/wopta/goworkspace/test v0.0.0-20230104102617-e971cb02bc2a
	github.com/wopta/goworkspace/form v0.0.0-20230104102617-e971cb02bc2a
	github.com/wopta/goworkspace/sellable v0.0.0
	github.com/wopta/goworkspace/question v0.0.0
	github.com/wopta/goworkspace/reserved v0.0.0
	github.com/wopta/goworkspace/companydata v0.0.0
	github.com/wopta/goworkspace/partnership v0.0.0
	"github.com/joho/godotenv" v1.5.1
)
