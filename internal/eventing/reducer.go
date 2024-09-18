package eventing

type Reducer interface {
	Reduce(events []Event)
}

// ChangeProducer produces changes for an aggregate.
type ChangeProducer interface {
	Changes() *AggregateChangeIntent
}

// Writer produces changes given a current view of the journal.
type Writer interface {
	JournalViewer
	ChangeProducer
}

type BaseWriter struct {
	aggregateType AggregateType
	aggregateID   AggregateID
	version       AggregateVersion
	events        []Event
	matcher       VersionMatcher
}

func NewBaseWriter(aggregateID AggregateID, aggregateType AggregateType, matcher VersionMatcher) *BaseWriter {
	return &BaseWriter{
		aggregateType: aggregateType,
		aggregateID:   aggregateID,
		matcher:       matcher,
	}
}

func (b *BaseWriter) Changes() *AggregateChangeIntent {
	return &AggregateChangeIntent{
		aggregateID:               b.aggregateID,
		aggregateType:             b.aggregateType,
		lastKnownAggregateVersion: b.version,
		events:                    b.events,
		versionMatcher:            b.matcher,
	}
}

func (b *BaseWriter) Reduce(events []*JournalEvent) {
	if len(events) == 0 {
		return
	}
	lastEvent := events[len(events)-1]
	b.version = lastEvent.AggregateVersion()
}

func (b *BaseWriter) Append(events ...Event) {
	b.events = append(b.events, events...)
}

func (b *BaseWriter) Aggregate() *Aggregate {
	return &Aggregate{
		AggregateID:   b.aggregateID,
		AggregateType: b.aggregateType,
		Version:       b.version,
	}
}

func (b *BaseWriter) Equals(another *BaseWriter) bool {
	if b == nil || another == nil {
		return false
	}
	return b.aggregateType == another.aggregateType &&
		b.aggregateID == another.aggregateID &&
		b.version == another.version &&
		len(b.events) == len(another.events)
}
