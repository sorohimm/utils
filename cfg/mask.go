package cfg

import (
	"fmt"
	"strings"
)

func FixedWidth(s string, maskBy string, limit int, unmasked int) string {
	if len(s) == 0 {
		return s
	}

	thirds := len(s) / 3
	rem := min(thirds, unmasked)
	tail := s[len(s)-rem:]

	mask := limit - rem

	maskedHead := strings.Repeat(maskBy, mask)
	masked := fmt.Sprintf("%v%v", maskedHead, tail)

	return masked
}

func min(x, y int) int {
	if x > y {
		return y
	}

	return x
}
