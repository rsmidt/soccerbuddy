package validation

import "time"

// ValidateStringRequiredWithLength validates a string that is required and has a minimum and maximum length.
func ValidateStringRequiredWithLength(s, name string, min, max int) FieldError {
	if s == "" {
		return NewFieldError(name, ErrRequired)
	} else if len(s) > max {
		return NewMaxLengthError(name, max)
	} else if len(s) < min {
		return NewMinLengthError(name, min)
	}
	return nil
}

// ValidateDateRequiredInRange validates a date that is required and has to be in a certain range.
func ValidateDateRequiredInRange(d time.Time, name string, min, max time.Time) FieldError {
	if d.IsZero() {
		return NewFieldError(name, ErrRequired)
	} else if d.Before(min) {
		return NewMinDateError(name, min)
	} else if d.After(max) {
		return NewMaxDateError(name, max)
	}
	return nil
}
