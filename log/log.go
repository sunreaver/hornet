package log

import (
	"fmt"
	"strings"
)

const header = "[hornet]"

// Logf Logf
func Logf(info ...interface{}) {
	fmt.Printf(makeFormat(len(info)), info...)
}

func makeFormat(count int) string {
	if count == 0 {
		return header
	}
	return fmt.Sprintf("%s %s\n", header, strings.Repeat(", %v", count)[2:])
}
