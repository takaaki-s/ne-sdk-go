package nextengine

// APIError is used for any errors that occur in the API
type APIError struct {
	commonResult
}

func (e *APIError) Error() string {
	return e.Message
}
