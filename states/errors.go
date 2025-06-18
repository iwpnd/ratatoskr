package states

import "fmt"

// StateError is a reusable custom error independent of the state it occurs.
type StateError struct {
	State States

	Err error
}

func (se *StateError) Error() string {
	return fmt.Sprintf("error in state %s (err: %s)", se.State.String(), se.Err.Error())
}
