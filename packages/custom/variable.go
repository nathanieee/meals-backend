package custom

import "strings"

type (
	V_STRING string
)

func (m *V_STRING) SuffixSpaceCheck() V_STRING {
	if !strings.HasSuffix(string(*m), " ") {
		*m += " "
	}

	return *m
}
