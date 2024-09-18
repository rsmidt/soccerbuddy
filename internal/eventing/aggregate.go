package eventing

import (
	"errors"
	"strconv"
)

type Aggregate struct {
	AggregateID   AggregateID
	AggregateType AggregateType
	Version       AggregateVersion
}

func (a *Aggregate) String() string {
	return string(a.AggregateType) + ":" + a.AggregateID.Deref() + ":" + strconv.FormatUint(uint64(a.Version), 10)
}

type AggregateChangeIntent struct {
	aggregateID               AggregateID
	aggregateType             AggregateType
	lastKnownAggregateVersion AggregateVersion
	events                    []Event
	versionMatcher            VersionMatcher
}

func (a *AggregateChangeIntent) AggregateID() AggregateID {
	return a.aggregateID
}

func (a *AggregateChangeIntent) AggregateType() AggregateType {
	return a.aggregateType
}

func (a *AggregateChangeIntent) LastKnownAggregateVersion() AggregateVersion {
	return a.lastKnownAggregateVersion
}

func (a *AggregateChangeIntent) Events() []Event {
	return a.events
}

func (a *AggregateChangeIntent) VersionMatches(actual AggregateVersion) bool {
	return a.versionMatcher.Matches(a.lastKnownAggregateVersion, actual)
}

// NewAggregateChangeIntent creates a new AggregateChangeIntent.
// The intent is to append the given events to the aggregate with the given ID and type.
// All events must be targeting the same aggregate.
func NewAggregateChangeIntent(aggregateID AggregateID, aggregateType AggregateType, lastKnownAggregateVersion AggregateVersion, events []Event, matcher VersionMatcher) (AggregateChangeIntent, error) {
	for _, event := range events {
		if event.AggregateID() != aggregateID {
			return AggregateChangeIntent{}, errors.New("event does not target the same aggregate")
		}
		if event.AggregateType() != aggregateType {
			return AggregateChangeIntent{}, errors.New("event does not target the same aggregate")
		}
	}

	return AggregateChangeIntent{
		aggregateID:               aggregateID,
		aggregateType:             aggregateType,
		lastKnownAggregateVersion: lastKnownAggregateVersion,
		events:                    events,
		versionMatcher:            matcher,
	}, nil
}
