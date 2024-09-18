package eventing

import (
	"fmt"
)

// UniqueConstraint represents a unique constraint on a field of an aggregate.
type UniqueConstraint struct {
	// ownerAggregateID is the ID of the aggregate that owns the unique constraint.
	ownerAggregateID AggregateID
	constrainedField string
	constrainedValue string
}

// NewUniqueConstraint creates a unique constraint for the given aggregate and field.
func NewUniqueConstraint(ownerAggregateID AggregateID, constrainedField, constraintValue string) UniqueConstraint {
	return UniqueConstraint{
		ownerAggregateID: ownerAggregateID,
		constrainedField: constrainedField,
		constrainedValue: constraintValue,
	}
}

// NewDeleteAllConstraint creates a unique constraint that removes all constraints for the given aggregate.
func NewDeleteAllConstraint(id AggregateID) UniqueConstraint {
	return UniqueConstraint{
		ownerAggregateID: id,
		constrainedField: "",
		constrainedValue: "",
	}
}

func (u *UniqueConstraint) OwnerAggregateID() AggregateID {
	return u.ownerAggregateID
}

func (u *UniqueConstraint) ConstrainedField() string {
	return u.constrainedField
}

func (u *UniqueConstraint) ConstrainedValue() string {
	return u.constrainedValue
}

type UniqueConstraintAdder interface {
	UniqueConstraintsToAdd() []UniqueConstraint
}

type UniqueConstraintRemover interface {
	UniqueConstraintsToRemove() []UniqueConstraint
}

func NewUniqueConstraintError(violated UniqueConstraint) error {
	return uniqueConstraintError{violated: violated}
}

type uniqueConstraintError struct {
	violated UniqueConstraint
}

func (u uniqueConstraintError) Error() string {
	return fmt.Sprintf("unique constraint violation on field %q with value %q", u.violated.ConstrainedField(), u.violated.ConstrainedValue())
}

func (u uniqueConstraintError) Violated() UniqueConstraint {
	return u.violated
}
