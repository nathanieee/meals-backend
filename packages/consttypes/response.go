package consttypes

type (
	ResponseStatusType string
)

const (
	RST_SUCCESS ResponseStatusType = "success"
	RST_FAIL    ResponseStatusType = "fail"
	RST_ERROR   ResponseStatusType = "error"
)

func (enum ResponseStatusType) String() string {
	return string(enum)
}
