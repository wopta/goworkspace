module github.com/wopta/goworkspace/lib

go 1.16

replace github.com/wopta/goworkspace/lib => ./lib

//require golang.org/x/oauth2/clientcredentials v0.0.0-20220622183110-fd043fe589d2

require (
	cloud.google.com/go/storage v1.22.1
	github.com/go-gota/gota v0.12.0
	github.com/ompluscator/dynamic-struct v1.3.0
	github.com/wopta/goworkspace/models v0.0.0-20220909121553-d232bcdeb3e0
	golang.org/x/crypto v0.0.0-20220829220503-c86fa9a7ed90
	golang.org/x/oauth2 v0.0.0-20220722155238-128564f6959c
)
