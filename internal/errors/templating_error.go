package errors

import (
    "fmt"
)

type TemplateError struct {
    Err error
}

func (te *TemplateError) Error() string {
    return fmt.Sprintf("Templating error: %s", te.Err.Error())
}

func NewTemplateError(err error) error {
    return &TemplateError{
        Err: err,
    }
}
