package toolchain

import "strings"

// SplitCommandTo Split Command and Adjust To.
func SplitCommandTo(raw string, setCommandStopper int) (splitCommandLen int, splitInfo []string) {
	rawSplit := strings.SplitN(raw, " ", setCommandStopper)
	return len(rawSplit), rawSplit
}
