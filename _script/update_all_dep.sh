#!/usr/bin/env zsh

location="callback"

cd document
go get github.com/wopta/goworkspace/lib
go get github.com/wopta/goworkspace/models
go get github.com/wopta/goworkspace/product
cd ..

cd rules
go get github.com/wopta/goworkspace/lib
go get github.com/wopta/goworkspace/models
go get github.com/wopta/goworkspace/quote
cd ..

cd broker
go get github.com/wopta/goworkspace/lib
go get github.com/wopta/goworkspace/models
go get github.com/wopta/goworkspace/document
go get github.com/wopta/goworkspace/mail
go get github.com/wopta/goworkspace/payment
cd ..

cd claim
go get github.com/wopta/goworkspace/lib
go get github.com/wopta/goworkspace/models
cd ..

cd enrich
go get github.com/wopta/goworkspace/lib
go get github.com/wopta/goworkspace/models
cd ..

cd payment
go get github.com/wopta/goworkspace/lib
go get github.com/wopta/goworkspace/models
cd ..

cd product
go get github.com/wopta/goworkspace/lib
go get github.com/wopta/goworkspace/models
cd ..

cd quote
go get github.com/wopta/goworkspace/lib
cd ..

cd "$location"
go get github.com/wopta/goworkspace/lib
go get github.com/wopta/goworkspace/models