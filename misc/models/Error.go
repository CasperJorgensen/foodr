package misc

type Error struct {
	ErrorText string    `json:"errorText,omitempty"`
	ErrorCode ErrorCode `json:"errorCode,omitempty"`
	IsError   bool      `json:"isError,omitempty"`
}

func NewError(text string, errorCode ErrorCode, isError bool) *Error {
	return &Error{
		ErrorText: text,
		ErrorCode: errorCode,
		IsError:   isError,
	}
}

//go:generate stringer -type=ErrorCode
type ErrorCode int

const (
	Default ErrorCode = iota
	EmailAlreadyExist
	UsernameAlreadyExist
	UsernameCannotBeNullOrEmpty
	FullnameCannotBeNullOrEmpty
	RecipeNotFound
)

func (b ErrorCode) MarshalText() ([]byte, error) {
	return []byte(b.String()), nil
}
