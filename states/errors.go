package states

import "fmt"

// StateError ...
type StateError struct {
	State States

	Err error
}

// Error ...
func (se *StateError) Error() string {
	return fmt.Sprintf("error in state %s (err: %s)", se.State.String(), se.Err.Error())
}
