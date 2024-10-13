package cfg

import (
	"fmt"
)

type SafeString string

func (t SafeString) MarshalJSON() ([]byte, error) {
	masked := FixedWidth(string(t), "*", 20, 5)

	return []byte(fmt.Sprintf("\"%s\"", masked)), nil
}

func (t SafeString) String() string {
	return string(t)
}
