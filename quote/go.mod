module github.com/wopta/goworkspace/quote

go 1.21

replace github.com/wopta/goworkspace/quote => ./

require (
	cloud.google.com/go v0.110.2
	github.com/GoogleCloudPlatform/functions-framework-go v1.7.3
	github.com/dustin/go-humanize v1.0.1
	github.com/go-gota/gota v0.12.0
	github.com/wopta/goworkspace/lib v1.0.99
	github.com/wopta/goworkspace/models v1.1.84
	github.com/wopta/goworkspace/network v1.0.36
	github.com/wopta/goworkspace/sellable v1.0.63
	github.com/xuri/excelize/v2 v2.8.0
	google.golang.org/api v0.122.0
	modernc.org/mathutil v1.5.0
)