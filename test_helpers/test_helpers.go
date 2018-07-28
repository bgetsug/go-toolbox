package test_helpers

import (
	"fmt"
	"os"
	"testing"
)

func Main(m *testing.M, suiteName string, initFunc func()) {
	initFunc()

	exit := m.Run()

	PrintLogs(suiteName)

	os.Exit(exit)
}

func PrintLogs(suiteName string) {
	if testing.Verbose() {
		fmt.Printf("\nLogs from all %s tests, not already captured\n", suiteName)
		fmt.Println("------------------------------------------------------------------------------")
		fmt.Print(SuiteLog.String())
		fmt.Println("------------------------------------------------------------------------------")
		fmt.Println()
	}
}
