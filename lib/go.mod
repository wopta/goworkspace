module github.com/wopta/goworkspace/lib

go 1.16

replace github.com/wopta/goworkspace/lib => ./lib

//require golang.org/x/oauth2/clientcredentials v0.0.0-20220622183110-fd043fe589d2

require (
	github.com/go-gota/gota v0.12.0
	golang.org/x/oauth2 v0.0.0-20220722155238-128564f6959c
)
