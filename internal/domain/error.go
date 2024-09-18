package domain

import (
	"fmt"
	"github.com/rsmidt/soccerbuddy/internal/eventing"
)

type InvalidAggregateStateError struct {
	Aggregate     *eventing.Aggregate
	ExpectedState int
	ActualState   int
}

func NewInvalidAggregateStateError(aggregate *eventing.Aggregate, expectedState, actualState int) *InvalidAggregateStateError {
	return &InvalidAggregateStateError{
		Aggregate:     aggregate,
		ExpectedState: expectedState,
		ActualState:   actualState,
	}
}

func (e InvalidAggregateStateError) Error() string {
	return fmt.Sprintf("invalid aggregate state: %s expected: %d actual: %d", e.Aggregate.String(), e.ExpectedState, e.ActualState)
}
