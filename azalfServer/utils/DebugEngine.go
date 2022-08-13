package utils

import "fmt"

// My own type for debugging the application
//    - Extensible for any type of data
//    - shows lines of code that are executed
//    - shows what type of debug message it is
type DebugStruct struct {
	FileLine  int32
	DebugType string
	Message   string
}

func CreateDebug(fileLine int32, debugType string, message string) DebugStruct {
	return DebugStruct{
		FileLine:  fileLine,
		DebugType: debugType,
		Message:   message,
	}
}

func (d *DebugStruct) String() string {
	if d.DebugType == "error" {
		// print in red
		return fmt.Sprintf("\x1b[31m%d: ERROR\t\x1b[0m%s\n", d.FileLine, d.Message)
	} else if d.DebugType == "warning" {
		return fmt.Sprintf("\x1b[33m%d: WARN\t\x1b[0m%s\n", d.FileLine, d.Message)
	} else if d.DebugType == "success" {
		// print in blue
		return fmt.Sprintf("\x1b[34m%d: SUCCESS\t\x1b[0m%s\n", d.FileLine, d.Message)
	} else {
		// print in gray
		return fmt.Sprintf("\x1b[37m%d: DEBUG\t\x1b[0m%s\n", d.FileLine, d.Message)
	}
}
