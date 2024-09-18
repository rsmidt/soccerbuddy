package eventing

// JournalInquirer creates a query for the journal.
type JournalInquirer interface {
	Query() JournalQuery
}

// JournalViewer takes a slice form the journal and reduces it to a view.
type JournalViewer interface {
	JournalInquirer

	Reduce(events []*JournalEvent)
}
