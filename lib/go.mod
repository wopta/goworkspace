module github.com/wopta/goworkspace/lib

go 1.16

replace github.com/wopta/goworkspace/lib => ./

require (
	cloud.google.com/go/firestore v1.8.0 // indirect
	cloud.google.com/go/storage v1.23.0
	firebase.google.com/go v3.13.0+incompatible
	github.com/go-gota/gota v0.12.0
	github.com/ompluscator/dynamic-struct v1.3.0
	github.com/rocketlaunchr/dataframe-go v0.0.0-20211025052708-a1030444159b
	golang.org/x/crypto v0.0.0-20220829220503-c86fa9a7ed90
	golang.org/x/oauth2 v0.0.0-20220909003341-f21342109be1
)
