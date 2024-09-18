package validation

import (
	"bytes"
	"fmt"
	"time"
)

// Errors represents a collection of field errors.
type Errors []FieldError

func (v Errors) Error() string {
	var buffer bytes.Buffer
	for i, fe := range v {
		if i > 0 {
			buffer.WriteString("\n")
		}
		buffer.WriteString(fe.Error())
	}
	return buffer.String()
}

func (v Errors) Is(err error) bool {
	if len(v) == 0 {
		return err == nil
	}
	if len(v) == 1 {
		return err.Error() == v[0].Error()
	}
	return err.Error() == v.Error()
}

// FieldError represents an error that occurred during validation of a field.
type FieldError interface {
	error

	Field() string
	Type() string
}

type genericFieldError struct {
	field string
	typ   string
}

func NewFieldError(field, typ string) FieldError {
	return &genericFieldError{
		field: field,
		typ:   typ,
	}
}

func (f genericFieldError) Error() string {
	return fmt.Sprintf("field %s: %s", f.field, f.typ)
}

func (f genericFieldError) Field() string {
	return f.field
}

func (f genericFieldError) Type() string {
	return f.typ
}

type minLengthError struct {
	field string
	min   int
}

func NewMinLengthError(field string, min int) FieldError {
	return &minLengthError{
		field: field,
		min:   min,
	}
}

func (m minLengthError) Error() string {
	return fmt.Sprintf("field %s: must be at least %d characters long", m.field, m.min)
}

func (m minLengthError) Field() string {
	return m.field
}

func (m minLengthError) Type() string {
	return ErrMinLength
}

type maxLengthError struct {
	field string
	max   int
}

func NewMaxLengthError(field string, max int) FieldError {
	return &maxLengthError{
		field: field,
		max:   max,
	}
}

func (m maxLengthError) Error() string {
	return fmt.Sprintf("field %s: must be at most %d characters long", m.field, m.max)
}

func (m maxLengthError) Field() string {
	return m.field
}

func (m maxLengthError) Type() string {
	return ErrMaxLength
}

type minDateError struct {
	field string
	min   time.Time
}

func NewMinDateError(field string, min time.Time) FieldError {
	return &minDateError{
		field: field,
		min:   min,
	}
}

func (m minDateError) Error() string {
	return fmt.Sprintf("field %s: must be at least after %s", m.field, m.min)
}

func (m minDateError) Field() string {
	return m.field
}

func (m minDateError) Type() string {
	return ErrMinDate
}

type maxDateError struct {
	field string
	max   time.Time
}

func NewMaxDateError(field string, max time.Time) FieldError {
	return &maxDateError{
		field: field,
		max:   max,
	}
}

func (m maxDateError) Error() string {
	return fmt.Sprintf("field %s: must be before %s", m.field, m.max)
}

func (m maxDateError) Field() string {
	return m.field
}

func (m maxDateError) Type() string {
	return ErrMaxDate
}

type existsError struct {
	field string
}

func NewExistsError(field string) FieldError {
	return &existsError{
		field: field,
	}
}

func (e existsError) Error() string {
	return fmt.Sprintf("field %s: already exists", e.field)
}

func (e existsError) Field() string {
	return e.field
}

func (e existsError) Type() string {
	return ErrAlreadyExists
}
