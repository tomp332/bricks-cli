package utils

import (
	"fmt"
	"github.com/fatih/color"
	"os"
)

// PrinterFunction defines the type for our printer functions
type PrinterFunction func(format string, args ...interface{})

// ErrorPrint prints error messages in red
func ErrorPrint(format string, args ...interface{}) {
	fmt.Println()
	color.New(color.FgRed).Printf("[-] "+format+"\n", args...)
	fmt.Println()
}

// FatalPrint prints error messages in red and exits the process
func FatalPrint(format string, args ...interface{}) {
	fmt.Println()
	color.New(color.FgRed).Printf("[-] "+format+"\n", args...)
	fmt.Println()
	os.Exit(1)
}

// SuccessPrint prints success messages in green
func SuccessPrint(format string, args ...interface{}) {
	fmt.Println()
	color.New(color.FgGreen).Printf("[+] "+format+"\n", args...)
	fmt.Println()
}

func InfoPrint(format string, args ...interface{}) {
	fmt.Println()
	color.New(color.FgCyan).Printf("[!] "+format+"\n", args...)
	fmt.Println()
}

// WarningPrint prints warning messages in yellow
func WarningPrint(format string, args ...interface{}) {
	fmt.Println()
	color.New(color.FgYellow).Printf("[*] "+format+"\n", args...)
	fmt.Println()
}
