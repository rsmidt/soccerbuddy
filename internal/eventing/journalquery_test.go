package eventing

import (
	"github.com/rsmidt/soccerbuddy/internal/core/idgen"
	"testing"
)

func TestJournalQueryBuilder(t *testing.T) {
	var builder JournalQueryBuilder
	idA := idgen.New[AggregateID]()
	idB := idgen.New[AggregateID]()
	query := builder.
		WithAggregate("test").
		AggregateID(idA).
		AggregateVersionAtLeast(1).
		Events("test").
		Finish().
		WithAggregate("test2").
		AggregateID(idB).
		Events().
		Finish().
		MustBuild()

	byType := query.AggQueriesByType()
	if len(byType) != 2 {
		t.Errorf("expected 2 types, got %d", len(byType))
	}
	if aggQuery, ok := byType["test"]; !ok {
		t.Errorf("expected type test")
	} else if aggQuery.version != 1 {
		t.Errorf("expected version 1, got %d", aggQuery.version)
	} else if len(aggQuery.events) != 1 {
		t.Errorf("expected 1 event, got %d", len(aggQuery.events))
	} else if aggQuery.events[0] != "test" {
		t.Errorf("expected event test, got %s", aggQuery.events[0])
	} else if aggQuery.id != idA {
		t.Errorf("expected id %s, got %s", idA, aggQuery.id)
	}

	if aggQuery, ok := byType["test2"]; !ok {
		t.Errorf("expected type test2")
	} else if aggQuery.version != 0 {
		t.Errorf("expected version 0, got %d", aggQuery.version)
	} else if len(aggQuery.events) != 0 {
		t.Errorf("expected 0 event, got %d", len(aggQuery.events))
	} else if aggQuery.id != idB {
		t.Errorf("expected id %s, got %s", idB, aggQuery.id)
	}
}
