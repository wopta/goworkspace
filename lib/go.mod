module github.com/wopta/goworkspace/lib

go 1.19

replace github.com/wopta/goworkspace/lib => ./

require (
	cloud.google.com/go/bigquery v1.44.0
	cloud.google.com/go/firestore v1.9.0
	cloud.google.com/go/storage v1.28.1
	firebase.google.com/go v3.13.0+incompatible
	github.com/dustin/go-humanize v0.0.0-20171111073723-bb3d318650d4
	github.com/go-gota/gota v0.12.0
	github.com/hyperjumptech/grule-rule-engine v1.12.0
	github.com/pkg/sftp v1.13.5
	github.com/rocketlaunchr/dataframe-go v0.0.0-20211025052708-a1030444159b
	github.com/xuri/excelize/v2 v2.7.0
	golang.org/x/crypto v0.5.0
	golang.org/x/oauth2 v0.0.0-20221014153046-6fdb5e3db783
	google.golang.org/api v0.103.0
)
