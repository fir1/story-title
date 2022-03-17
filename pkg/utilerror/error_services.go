package utilerror

// ErrArgument is used when an argument that came from the client is erroneous.
// Example: an attempt to log in with the empty password.
type ErrArgument struct {
	Wrapped error
}

func (e ErrArgument) Error() string {
	return "invalid argument"
}

func (e ErrArgument) Unwrap() error {
	return e.Wrapped
}

// ErrForbidden is used when a person is forbidden to proceed with an action.
type ErrForbidden struct {
	Wrapped error
}

func (e ErrForbidden) Error() string {
	return "forbidden"
}

func (e ErrForbidden) Unwrap() error {
	return e.Wrapped
}

// ErrUnauthorised is used when a person has presented an expired or non-existent session.
type ErrUnauthorised struct {
	Wrapped error
}

func (e ErrUnauthorised) Error() string {
	return "forbidden"
}

func (e ErrUnauthorised) Unwrap() error {
	return e.Wrapped
}

// ErrPrecondition is used when the system is in a state that prevents the action that
// is not involved with authentication.
type ErrPrecondition struct {
	Wrapped error
}

func (e ErrPrecondition) Error() string {
	return "failed precondition"
}

func (e ErrPrecondition) Unwrap() error {
	return e.Wrapped
}
