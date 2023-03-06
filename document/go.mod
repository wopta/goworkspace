module github.com/wopta/goworkspace/document

go 1.16

replace github.com/wopta/goworkspace/document => ./

require (
	firebase.google.com/go/v4 v4.10.0 // indirect
	github.com/GoogleCloudPlatform/functions-framework-go v1.6.1
	github.com/dustin/go-humanize v1.0.1
	github.com/johnfercher/maroto v0.38.0
	github.com/ruudk/golang-pdf417 v0.0.0-20201230142125-a7e3863a1245 // indirect
	github.com/wopta/goworkspace/lib v0.0.0-20230306085301-1c04a6a66d81
	github.com/wopta/goworkspace/models v0.0.0-20230306085301-1c04a6a66d81
	github.com/wopta/goworkspace/product v0.0.0-20230306085301-1c04a6a66d81
)
