module github.com/wopta/goworkspace/rules

go 1.16

replace github.com/wopta/goworkspace/rules => ./rules

require github.com/hyperjumptech/grule-rule-engine v1.11.0

require github.com/jinzhu/copier v0.3.5

require github.com/ompluscator/dynamic-struct v1.3.0

require github.com/leebenson/conform v1.2.2

require cloud.google.com/go/storage v1.25.0

require (
	github.com/golang/protobuf v1.5.2 // indirect
	golang.org/x/net v0.0.0-20220809184613-07c6da5e1ced // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)
