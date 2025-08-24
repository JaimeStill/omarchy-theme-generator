package tests

import "fmt"

func CheckMark(condition bool) string {
	if condition {
		return "✓"
	}
	return "✗"
}

func TestSection(name string) {
	fmt.Printf("\n=== %s ===\n", name)
}

func TestResult(name string, pass bool, reason string) {
	fmt.Printf("%s: %s - %s\n", name, CheckMark(pass), reason)
}

func ShowTransformation(before, operation, after string) {
	fmt.Printf("  %s → [%s] → %s\n", before, operation, after)
}

func CompareValues(name string, expected, actual any, pass bool) {
	fmt.Printf("  %s: expected=%v, actual=%v %s\n", name, expected, actual, CheckMark(pass))
}
