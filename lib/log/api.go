package log

import (
	"encoding/json"
	"fmt"
)

// Append the prefix, ex: [prefix1] -> [prefix1|prefix2]
// Remember to use PopPrefix to remove it eventually
func AddPrefix(prefix string) {
	Log().AddPrefix(prefix)
}

// Remove the younger prefix, ex: [prefix1|prefix2] -> [prefix1]
func PopPrefix() {
	Log().PopPrefix()
}

// Remove all prefixs, ex: [prefix1|prefix2] -> <None>
func ResetPrefix() {
	Log().ResetPrefix()
}

// Log a formatted message with severity 'DEFAULT'
func Printf(format string, a ...any) {
	Log().Printf(format, a...)
}

// Log a message with severity 'DEFAULT'
func Println(message ...any) {
	Log().Println(fmt.Sprint(message...))
}

// Log a formatted message with severity 'INFO'
func InfoF(format string, a ...any) {
	Log().InfoF(format, a...)
}

// Log a formatted message with severity 'WARNING'
func WarningF(format string, a ...any) {
	Log().WarningF(format, a...)
}

// Log a error, with struct : 'Error: <err>'
func Error(err error) {
	Log().Error(err)
}

// Log a formatted message with severity 'ERROR'
func ErrorF(format string, a ...any) {
	Log().ErrorF(format, a...)
}

// Log a struct
func PrintStruct(message string, object any) {
	bytes, err := json.Marshal(object)
	if err != nil {
		Error(err)
		return
	}
	Printf(message+": %v", string(bytes))
}
