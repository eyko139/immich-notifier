package errors

import (
    "fmt"
)

type NoUserInSessionError struct {
    Message string
    Session string
}

func (e *NoUserInSessionError) Error() string {
    return fmt.Sprintf("No found in Session: %s; err: %s", e.Session, e.Message)
}
