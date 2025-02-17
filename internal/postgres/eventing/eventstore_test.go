package eventing

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/rsmidt/soccerbuddy/internal/core"
	"github.com/rsmidt/soccerbuddy/internal/core/idgen"
	"github.com/rsmidt/soccerbuddy/internal/eventing"
	"github.com/rsmidt/soccerbuddy/internal/postgres"
)

var _ eventing.JournalEventMapper = (*testRegistry)(nil)

type testRegistry struct{}

func (*testRegistry) MapFrom(
	aggregateID eventing.AggregateID,
	aggregateType eventing.AggregateType,
	eventVersion eventing.EventVersion,
	eventType eventing.EventType,
	eventID eventing.EventID,
	aggregateVersion eventing.AggregateVersion,
	journalPosition eventing.JournalPosition,
	insertedAt time.Time,
	payload []byte,
) (*eventing.JournalEvent, error) {
	base := eventing.NewEventBase(aggregateID, aggregateType, eventVersion, eventType)
	combinedID := fmt.Sprintf("%s::%s%s", aggregateType, eventType, eventVersion)

	switch combinedID {
	case "test::TestEventv1":
		event := &testEvent{
			EventBase: base,
		}
		if err := json.Unmarshal(payload, event); err != nil {
			return nil, fmt.Errorf("failed to unmarshal event: %w", err)
		}
		return eventing.NewJournalEvent(event, eventID, aggregateVersion, journalPosition, insertedAt), nil
	case "test::TestEventUniqueConstraintv1":
		event := &testEventUniqueConstraint{
			EventBase: base,
		}
		if err := json.Unmarshal(payload, event); err != nil {
			return nil, fmt.Errorf("failed to unmarshal event: %w", err)
		}
		return eventing.NewJournalEvent(event, eventID, aggregateVersion, journalPosition, insertedAt), nil
	default:
		return nil, errors.New("event not registered")
	}
}

type stubCrypto struct{}

func (s *stubCrypto) EncryptEvents(ctx context.Context, events []eventing.Event) error {
	return nil
}

func (s *stubCrypto) DecryptEvents(ctx context.Context, events []eventing.Event) error {
	return nil
}

func newTestEvent(id string) *testEvent {
	return &testEvent{
		EventBase: eventing.NewEventBase(eventing.AggregateID(id), "test", "v1", "TestEvent"),
		ID:        id,
		FieldA:    "test",
	}
}

type testEvent struct {
	*eventing.EventBase

	ID     string `json:"id"`
	FieldA string `json:"field_a"`
}

func (t *testEvent) IsShredded() bool {
	return false
}

func newTestEventUniqueConstraint(id string) *testEventUniqueConstraint {
	return &testEventUniqueConstraint{
		EventBase: eventing.NewEventBase(eventing.AggregateID(id), "test", "v1", "TestEventUniqueConstraint"),
		ID:        id,
		FieldA:    "test",
	}
}

type testEventUniqueConstraint struct {
	*eventing.EventBase

	ID     string `json:"id"`
	FieldA string `json:"field_a"`
}

func (t *testEventUniqueConstraint) IsShredded() bool {
	return false
}

func (t *testEventUniqueConstraint) UniqueConstraintsToAdd() []eventing.UniqueConstraint {
	return []eventing.UniqueConstraint{
		eventing.NewUniqueConstraint(t.AggregateID(), "field_a", t.FieldA),
	}
}

func Test_pgEventStore_Append(t *testing.T) {
	type fields struct {
	}
	type args struct {
		ctx     context.Context
		intents []eventing.AggregateChangeIntent
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		verify func(t *testing.T, events []*eventing.JournalEvent, err error)
	}{
		{
			name: "appends an intent with just one event",
			args: args{
				ctx: context.Background(),
				intents: []eventing.AggregateChangeIntent{
					core.Must2(eventing.NewAggregateChangeIntent(
						"9e0563c7-9b1d-47d7-b120-d6aa2f1db7e1",
						"test",
						0,
						[]eventing.Event{newTestEvent("9e0563c7-9b1d-47d7-b120-d6aa2f1db7e1")},
						eventing.VersionMatcherAlways,
					)),
				},
			},
			verify: func(t *testing.T, events []*eventing.JournalEvent, err error) {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if len(events) != 1 {
					t.Fatalf("expected 1 event, got %d", len(events))
				}

				event := events[0]
				if event.AggregateVersion() != 1 {
					t.Errorf("expected aggregate version 1, got %d", event.AggregateVersion())
				}
				if event.AggregateID() != "9e0563c7-9b1d-47d7-b120-d6aa2f1db7e1" {
					t.Errorf("expected aggregate ID 9e0563c7-9b1d-47d7-b120-d6aa2f1db7e1, got %s", event.AggregateID())
				}

				testEvent, ok := event.Event.(*testEvent)
				if !ok {
					t.Fatalf("expected *testEvent, got %T", event.Event)
				}
				if testEvent.FieldA != "test" {
					t.Errorf("expected FieldA 'test', got %s", testEvent.FieldA)
				}
			},
		}, {
			name: "fails an intent which violates unique constraint",
			args: args{
				ctx: context.Background(),
				intents: []eventing.AggregateChangeIntent{
					core.Must2(eventing.NewAggregateChangeIntent(
						"9e0563c7-9b1d-47d7-b120-d6aa2f1db7e1",
						"test",
						0,
						[]eventing.Event{
							newTestEventUniqueConstraint("9e0563c7-9b1d-47d7-b120-d6aa2f1db7e1"),
							newTestEventUniqueConstraint("9e0563c7-9b1d-47d7-b120-d6aa2f1db7e1"),
						},
						eventing.VersionMatcherAlways,
					)),
				},
			},
			verify: func(t *testing.T, events []*eventing.JournalEvent, err error) {
				expectedErr := eventing.NewUniqueConstraintError(
					eventing.NewUniqueConstraint(
						"9e0563c7-9b1d-47d7-b120-d6aa2f1db7e1",
						"field_a",
						"test"))
				if !errors.Is(err, expectedErr) {
					t.Errorf("expected error %v, got %v", expectedErr, err)
				}
				if len(events) != 0 {
					t.Errorf("expected 0 events, got %d", len(events))
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			pool, cleanup := postgres.GetTestPool()
			t.Cleanup(cleanup)

			p := NewEventStore(postgres.GetTestLogger(), pool, &testRegistry{}, &stubCrypto{})
			events, err := p.Append(tt.args.ctx, tt.args.intents...)
			tt.verify(t, events, err)
		})
	}

	t.Run("appends multiple events with correct sequence", func(t *testing.T) {
		t.Parallel()

		pool, cleanup := postgres.GetTestPool()
		t.Cleanup(cleanup)

		es := NewEventStore(postgres.GetTestLogger(), pool, &testRegistry{}, &stubCrypto{})

		aggregateID := idgen.New[eventing.AggregateID]()
		aggregateType := eventing.AggregateType("test")

		// Create three events
		events := []eventing.Event{
			newTestEvent(string(aggregateID)),
			newTestEvent(string(aggregateID)),
			newTestEvent(string(aggregateID)),
		}

		journalEvents, err := es.Append(context.Background(), core.Must2(eventing.NewAggregateChangeIntent(
			aggregateID,
			aggregateType,
			0,
			events,
			eventing.VersionMatcherAlways,
		)))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(journalEvents) != 3 {
			t.Fatalf("expected 3 events, got %d", len(journalEvents))
		}

		// Verify sequence
		for i, event := range journalEvents {
			expectedVersion := eventing.AggregateVersion(i + 1)
			if event.AggregateVersion() != expectedVersion {
				t.Errorf("event %d: expected version %d, got %d", i, expectedVersion, event.AggregateVersion())
			}
			if event.AggregateID() != aggregateID {
				t.Errorf("event %d: expected aggregate ID %s, got %s", i, aggregateID, event.AggregateID())
			}
			if event.AggregateType() != aggregateType {
				t.Errorf("event %d: expected aggregate type %s, got %s", i, aggregateType, event.AggregateType())
			}

			// Verify the journal position is increasing
			if i > 0 && event.JournalPosition().Deref().Compare(journalEvents[i-1].JournalPosition().Deref()) < 1 {
				t.Errorf("event %d: journal position should be strictly increasing, got %v <= %v",
					i, event.JournalPosition(), journalEvents[i-1].JournalPosition())
			}
		}
	})

	t.Run("fails if sequences do not match", func(t *testing.T) {
		t.Parallel()

		pool, cleanup := postgres.GetTestPool()
		t.Cleanup(cleanup)

		es := NewEventStore(postgres.GetTestLogger(), pool, &testRegistry{}, &stubCrypto{})

		aggregateID := idgen.New[eventing.AggregateID]()
		aggregateType := eventing.AggregateType("test")

		// First append
		events, err := es.Append(context.Background(), core.Must2(eventing.NewAggregateChangeIntent(
			aggregateID,
			aggregateType,
			0,
			[]eventing.Event{newTestEvent(string(aggregateID))},
			eventing.VersionMatcherAlways,
		)))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(events) != 1 {
			t.Fatalf("expected 1 event, got %d", len(events))
		}
		if events[0].AggregateVersion() != 1 {
			t.Errorf("expected version 1, got %d", events[0].AggregateVersion())
		}

		// Second append with wrong version matcher
		events, err = es.Append(context.Background(), core.Must2(eventing.NewAggregateChangeIntent(
			aggregateID,
			aggregateType,
			0,
			[]eventing.Event{newTestEvent(string(aggregateID))},
			eventing.VersionMatcherExact,
		)))
		if !errors.Is(err, eventing.ErrVersionMismatch) {
			t.Errorf("expected ErrVersionMismatch, got %v", err)
		}
		if len(events) != 0 {
			t.Errorf("expected 0 events, got %d", len(events))
		}
	})
}

func TestMarshalsEventWithBaseEvent(t *testing.T) {
	t.Parallel()

	event := newTestEvent(idgen.NewString())
	result, err := json.Marshal(event)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expected := fmt.Sprintf(`{"id":"%s","field_a":"%s"}`, event.ID, event.FieldA)
	if string(result) != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}

	unmarshalled := &testEvent{
		EventBase: eventing.NewEventBase(event.AggregateID(), event.AggregateType(), event.EventVersion(), event.EventType()),
	}
	err = json.Unmarshal(result, unmarshalled)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	// Check that the base event was not overridden.
	if unmarshalled.AggregateID() != event.AggregateID() {
		t.Errorf("expected %s, got %s", event.AggregateID(), unmarshalled.AggregateID())
	}
	if unmarshalled.AggregateType() != event.AggregateType() {
		t.Errorf("expected %s, got %s", event.AggregateType(), unmarshalled.AggregateType())
	}
	if unmarshalled.EventVersion() != event.EventVersion() {
		t.Errorf("expected %s, got %s", event.EventVersion(), unmarshalled.EventVersion())
	}
	if unmarshalled.EventType() != event.EventType() {
		t.Errorf("expected %s, got %s", event.EventType(), unmarshalled.EventType())
	}
	// Check that the event-specific fields were unmarshalled.
	if unmarshalled.ID != event.ID {
		t.Errorf("expected %s, got %s", event.ID, unmarshalled.ID)
	}
	if unmarshalled.FieldA != event.FieldA {
		t.Errorf("expected %s, got %s", event.FieldA, unmarshalled.FieldA)
	}
}

func TestPgEventStore_Query(t *testing.T) {
	t.Run("queries events", func(t *testing.T) {
		t.Parallel()

		pool, cleanup := postgres.GetTestPool()
		t.Cleanup(cleanup)

		es := NewEventStore(postgres.GetTestLogger(), pool, &testRegistry{}, &stubCrypto{})
		eventA := newTestEvent(idgen.NewString())
		core.Must2(es.Append(context.Background(), core.Must2(eventing.NewAggregateChangeIntent(
			eventA.AggregateID(),
			eventA.AggregateType(),
			0,
			[]eventing.Event{eventA},
			eventing.VersionMatcherAlways,
		))))
		eventB := newTestEvent(idgen.NewString())
		core.Must2(es.Append(context.Background(), core.Must2(eventing.NewAggregateChangeIntent(
			eventB.AggregateID(),
			eventB.AggregateType(),
			0,
			[]eventing.Event{eventB},
			eventing.VersionMatcherAlways,
		))))

		var queryBuilder eventing.JournalQueryBuilder
		query := queryBuilder.
			WithAggregate("test").
			AggregateID(eventA.AggregateID()).
			AggregateVersionAtLeast(0).
			Events("TestEvent", "TestEvent2").
			Finish().MustBuild()

		events, err := es.Query(context.Background(), query)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(events) != 1 {
			t.Fatalf("expected 1 event, got %d", len(events))
		}
		result := events[0]
		if result.AggregateID() != eventA.AggregateID() {
			t.Fatalf("expected %s, got %s", eventA.AggregateID(), events[0].AggregateID())
		}
		if result.AggregateType() != eventA.AggregateType() {
			t.Fatalf("expected %s, got %s", eventA.AggregateType(), events[0].AggregateType())
		}
		if result.EventVersion() != eventA.EventVersion() {
			t.Fatalf("expected %s, got %s", eventA.EventVersion(), events[0].EventVersion())
		}
		if result.EventType() != eventA.EventType() {
			t.Fatalf("expected %s, got %s", eventA.EventType(), events[0].EventType())
		}
		casted, ok := result.Event.(*testEvent)
		if !ok {
			t.Fatalf("expected testEvent, got %T", result.Event)
		}
		if casted.ID != eventA.ID {
			t.Fatalf("expected %s, got %s", eventA.ID, casted.ID)
		}
		if casted.FieldA != eventA.FieldA {
			t.Fatalf("expected %s, got %s", eventA.FieldA, casted.FieldA)
		}
	})
}
