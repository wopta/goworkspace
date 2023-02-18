

Rem .\push.bat  v0.0.2 init_rules
set location="callback"
set v="0.6.0.dev"
set m="major allineament"
cd ..
git add .
git commit -m %m%
git push work master 
git push google master  
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
cd %location%
go get github.com/wopta/goworkspace/lib
go get github.com/wopta/goworkspace/models
cd ..
cd form
go get github.com/wopta/goworkspace/lib