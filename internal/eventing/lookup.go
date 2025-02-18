package eventing

import "errors"

var (
	ErrOwnerNotFound = errors.New("owner not found")
	ErrValueNotFound = errors.New("value not found")
)

type (
	LookupFieldName  string
	LookupFieldValue string
	LookupMap        map[LookupFieldName]LookupFieldValue
)

type LookupOpts struct {
	AggregateID   AggregateID
	AggregateType AggregateType
	FieldName     LookupFieldName
	FieldValue    LookupFieldValue
}

// LookupProvider is an interface that can be implemented by an [Event] to provide lookup values.
type LookupProvider interface {
	LookupValues() LookupMap
}

// LookupRemover is an interface that can be implemented by an [Event] to provide lookup values to remove.
type LookupRemover interface {
	LookupRemoves() []LookupFieldName
}

func (f *LookupFieldValue) Deref() string {
	return string(*f)
}
