package main

type UserAlreadyExistsError struct {
	message string
}

func (e *UserAlreadyExistsError) Error() string {
	return e.message
}
