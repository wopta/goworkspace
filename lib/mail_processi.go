package lib

import (
	"os"
	"strings"
)

func GetMailProcessi(sub string) string {
	mail := os.Getenv("MAIL_PROCESSI")
	pieces := strings.Split(mail, "@")
	return pieces[0] + "+" + sub + "@" + pieces[1]
}
