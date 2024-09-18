package domain

// Operator is a value object that specifies who operated an event.
type Operator struct {
	ActorID AccountID `json:"actor_id"`

	// OnBehalfOf is only set if the operator is acting on behalf of another person.
	// E.g., the system admin is not associated with a person.
	OnBehalfOf *PersonID `json:"on_behalf_of"`
}

// NewOperator creates a new Operator.
func NewOperator(actorID AccountID, onBehalfOf *PersonID) Operator {
	return Operator{
		ActorID:    actorID,
		OnBehalfOf: onBehalfOf,
	}
}
