package errors

type APIError struct {
	HTTPErrCode   int
	Description   string
	InternalError error
}

func NewAPIError(httpeErrCode int, description string, err error) *APIError {
	return &APIError{
		HTTPErrCode:   httpeErrCode,
		Description:   description,
		InternalError: err,
	}
}

func (e *APIError) Error() string {
	return e.Description
}
