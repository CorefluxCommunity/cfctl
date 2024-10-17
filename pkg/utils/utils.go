package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/logrusorgru/aurora"
)

const (
	printKVPadWidth     int = 30
	printHeaderPadWidth int = 10
)

// PrintKV will bold print the key followed by padding to the specified
// total width, then the value.
func PrintKV(key string, value string) {
	pad := strings.Repeat(".", printKVPadWidth-len(key)-2)
	label := aurora.Bold(fmt.Sprintf("%s %s:", key, pad))
	fmt.Println(label, value)
}

// PrintKV will bold print the key followed by padding to the specified
// total width, then the first value in the slice. If additional values
// are present in the slice they will be displayed on new line indented
// to match the previous value.
func PrintKVSlice(key string, values []string) {
	for i, value := range values {
		if i == 0 {
			pad := strings.Repeat(".", printKVPadWidth-len(key)-2)
			label := aurora.Bold(fmt.Sprintf("%s %s:", key, pad))
			fmt.Println(label, value)
		} else {
			pad := strings.Repeat(" ", printKVPadWidth)
			fmt.Println(pad, value)
		}
	}
}

// PrintHeader will print a bolded header label.
func PrintHeader(label string) {
	pad := strings.Repeat("=", printHeaderPadWidth)
	fmt.Println(aurora.Bold(fmt.Sprintf("%s %s %s", pad, label, pad)))
}

// PrintInfo will print a formatted info message to stdout.
func PrintInfo(msg string) {
	fmt.Println(aurora.Blue(aurora.Bold("[info]   ")), msg)
}

// PrintSuccess will print a formatted success message to stdout.
func PrintSuccess(msg string) {
	fmt.Println(aurora.Green(aurora.Bold("[success]")), msg)
}

// PrintWarning will print a formatted warning message to stdout.
func PrintWarning(msg string) {
	fmt.Println(aurora.Yellow(aurora.Bold("[warning]")), msg)
}

// PrintError will print a formatted error message to stdout.
func PrintError(msg string) {
	fmt.Println(aurora.Red(aurora.Bold("[error]  ")), msg)
}

// PrintFatal will print a formatted error message to stdout and exit with
// the provided status.
func PrintFatal(msg string, code int) {
	fmt.Println(aurora.Red(aurora.Bold("[fatal]  ")), msg)
	os.Exit(code)
}
