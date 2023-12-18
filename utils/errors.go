package utils

type NotFoundError struct {
	Trace error
}

func (e *NotFoundError) Error() string {
	return e.Trace.Error()
}

type ConflictError struct {
	Trace error
}

func (e *ConflictError) Error() string {
	return e.Trace.Error()
}
