package errors

type ErrorType struct {
	t string
}

var (
	ErrorTypeUnknown        = ErrorType{"unknown"}
	ErrorTypeAuthorization  = ErrorType{"authorization"}
	ErrorTypeIncorrectInput = ErrorType{"incorrect-input"}
	ErrorTypeNotFound       = ErrorType{"not-found"}
)

type CustomError struct {
	error     string
	message   string
	errorType ErrorType
}

func (s CustomError) Error() string {
	return s.error
}

func (s CustomError) Message() string {
	return s.message
}

func (s CustomError) ErrorType() ErrorType {
	return s.errorType
}

func NewCustomError(error string, message string) CustomError {
	return CustomError{
		error:     error,
		message:   message,
		errorType: ErrorTypeUnknown,
	}
}

func NewAuthorizationError(error string, message string) CustomError {
	return CustomError{
		error:     error,
		message:   message,
		errorType: ErrorTypeAuthorization,
	}
}

func NewNotFoundError(error string, message string) CustomError {
	return CustomError{
		error:     error,
		message:   message,
		errorType: ErrorTypeNotFound,
	}
}

func NewIncorrectInputError(error string, message string) CustomError {
	return CustomError{
		error:     error,
		message:   message,
		errorType: ErrorTypeIncorrectInput,
	}
}
