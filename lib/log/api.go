package log

import (
	"encoding/json"
	"fmt"
)

// Append the prefix, ex: [prefix1] -> [prefix1|prefix2]
// Remember to use PopPrefix to remove it eventually
func AddPrefix(prefix string) {
	Log().addPrefix(prefix)
}

// Remove the younger prefix, ex: [prefix1|prefix2] -> [prefix1]
func PopPrefix() {
	Log().popPrefix()
}

// Remove all prefixs, ex: [prefix1|prefix2] -> <None>
func ResetPrefix() {
	Log().resetPrefix()
}

// Log a formatted message with severity 'DEFAULT'
func Printf(format string, a ...any) {
	Log().printf(format, a...)
}

// Log a message with severity 'DEFAULT'
func Println(message ...any) {
	Log().println(fmt.Sprint(message...))
}

// Log a formatted message with severity 'INFO'
func InfoF(format string, a ...any) {
	Log().infoF(format, a...)
}

// Log a formatted message with severity 'WARNING'
func WarningF(format string, a ...any) {
	Log().warningF(format, a...)
}

// Log a error, with struct : 'Error: <err>'
func Error(err error) {
	Log().error(err)
}

// Log a formatted message with severity 'ERROR'
func ErrorF(format string, a ...any) {
	Log().errorF(format, a...)
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
