package main

import (
	"flag"
	"fmt"

	"dev.wopta.it/cloudbuild/scripts/tag_functions"
	"dev.wopta.it/cloudbuild/scripts/tag_modules"
)

var (
	script = flag.String("script", "", "the script to execute")
)

func main() {
	flag.Parse()

	if *script == "" {
		panic("choose script to run")
	}

	switch *script {
	case "tag_modules":
		tag_modules.Exec()
	case "tag_functions":
		tag_functions.Exec()
	default:
		panic(fmt.Sprintf("unknown script '%s'", *script))
	}
	fmt.Println("Script completed!")
}
