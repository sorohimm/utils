package cfg

import (
	"bytes"
	"strings"
)

type Error struct {
	errors []error
}

func (e Error) Error() string {
	buff := bytes.NewBufferString("")

	for i := 0; i < len(e.errors); i++ {
		err := e.errors[i]
		buff.WriteString(err.Error())
		buff.WriteString("\n")
	}

	return strings.TrimSpace(buff.String())
}
