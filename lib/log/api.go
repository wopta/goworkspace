package log

import "fmt"

func AddPrefix(prefix string) {
	Log().AddPrefix(prefix)
}

func PopPrefix() {
	Log().PopPrefix()
}

func ResetPrefix() {
	Log().ResetPrefix()
}
func Printf(format string, a ...any) {
	Log().Printf(format, a...)
}
func Println(message ...any) {
	Log().Println(fmt.Sprint(message...))
}
func InfoF(format string, a ...any) {
	Log().InfoF(format, a...)
}
func WarningF(format string, a ...any) {
	Log().WarningF(format, a...)
}
func Error(err error) {
	Log().Error(err)
}
func ErrorF(format string, a ...any) {
	Log().ErrorF(format, a...)
}
