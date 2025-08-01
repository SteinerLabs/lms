package web

type Error struct {
	Status  int    `json:"-"`
	Message string `json:"error"`
}

func NewError(code int, message string) *Error {
	return &Error{
		Status:  code,
		Message: message,
	}
}

func (e *Error) Error() string {
	return e.Message
}
