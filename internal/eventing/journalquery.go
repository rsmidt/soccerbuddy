package eventing

type JournalQuery struct {
	byType               map[AggregateType]AggregateQuery
	journalPositionAfter *JournalPosition
}

func (q *JournalQuery) AggQueriesByType() map[AggregateType]AggregateQuery {
	return q.byType
}

func (q *JournalQuery) JournalPositionAfter() *JournalPosition {
	return q.journalPositionAfter
}

type AggregateQuery struct {
	id      AggregateID
	version AggregateVersion
	events  []EventType
}

func (q *AggregateQuery) ID() AggregateID {
	return q.id
}

func (q *AggregateQuery) Version() AggregateVersion {
	return q.version
}

func (q *AggregateQuery) Events() []EventType {
	return q.events
}

type JournalQueryBuilder struct {
	byType               map[AggregateType]AggregateQuery
	journalPositionAfter *JournalPosition
}

// NewJournalQueryBuilderFrom creates a new JournalQueryBuilder using an existing JournalQuery.
func NewJournalQueryBuilderFrom(query JournalQuery) *JournalQueryBuilder {
	return &JournalQueryBuilder{
		byType:               query.byType,
		journalPositionAfter: query.journalPositionAfter,
	}
}

func (d *JournalQueryBuilder) WithAggregate(typ AggregateType) *AggregateQueryBuilder {
	return &AggregateQueryBuilder{
		jq:  d,
		typ: typ,
	}
}

func (d *JournalQueryBuilder) WithJournalPositionAfter(pos JournalPosition) *JournalQueryBuilder {
	d.journalPositionAfter = &pos
	return d
}

func (d *JournalQueryBuilder) MustBuild() JournalQuery {
	return JournalQuery{
		byType:               d.byType,
		journalPositionAfter: d.journalPositionAfter,
	}
}

type AggregateQueryBuilder struct {
	jq *JournalQueryBuilder

	id      AggregateID
	version AggregateVersion
	typ     AggregateType
	events  []EventType
}

func (d *AggregateQueryBuilder) AggregateID(id AggregateID) *AggregateQueryBuilder {
	d.id = id
	return d
}

func (d *AggregateQueryBuilder) AggregateVersionAtLeast(version AggregateVersion) *AggregateQueryBuilder {
	d.version = version
	return d
}

func (d *AggregateQueryBuilder) Events(eventTypes ...EventType) *AggregateQueryBuilder {
	d.events = eventTypes
	return d
}

func (d *AggregateQueryBuilder) Finish() *JournalQueryBuilder {
	if d.jq.byType == nil {
		d.jq.byType = make(map[AggregateType]AggregateQuery)
	}
	d.jq.byType[d.typ] = AggregateQuery{
		id:      d.id,
		version: d.version,
		events:  d.events,
	}
	return d.jq
}
