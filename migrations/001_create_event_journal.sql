-- Enable citext extension if not already enabled
CREATE EXTENSION IF NOT EXISTS "citext";

-- Create a function that returns the current timestamp.
CREATE OR REPLACE FUNCTION global_position()
    RETURNS numeric
AS
$$
BEGIN
    RETURN (SELECT EXTRACT(EPOCH FROM clock_timestamp()));
END;
$$ LANGUAGE plpgsql;

-- Create or replace the trigger function that notifies the event store channel.
CREATE OR REPLACE FUNCTION notify_event_store()
    RETURNS trigger AS $$
DECLARE
    channel_name TEXT;
BEGIN
    -- Set the channel name based on the aggregate_type.
    channel_name := 'event_store_' || NEW.aggregate_type;

    -- Send the notification to the grouped channel
    PERFORM pg_notify(channel_name, NEW.event_type);

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


-- Create the main event journal table.
CREATE TABLE event_journal
(
    id                TEXT    NOT NULL UNIQUE,

    aggregate_id      TEXT    NOT NULL,
    aggregate_type    TEXT    NOT NULL,
    aggregate_version BIGINT  NOT NULL,

    global_position   NUMERIC NOT NULL         DEFAULT global_position(),

    event_type        TEXT    NOT NULL,
    event_version     TEXT    NOT NULL,
    payload           JSONB   NOT NULL,
    metadata          JSONB,

    created_at        TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (aggregate_version, aggregate_id, aggregate_type)
);

-- Create index for sequence-based querying.
CREATE INDEX idx_event_journal_current_version ON event_journal (aggregate_version desc,
                                                                 aggregate_id asc, aggregate_type
                                                                 asc);

-- Create index for aggregate ID-based querying.
CREATE INDEX idx_event_journal_aggregate ON event_journal (aggregate_id, aggregate_type);

-- Create index for projection querying.
CREATE INDEX idx_event_journal_projection ON event_journal (aggregate_type, event_type, global_position);

-- Create the trigger to notify the event store channel.
CREATE TRIGGER event_store_notify_trigger
    AFTER INSERT ON event_journal  -- Replace with your actual table name
    FOR EACH ROW
EXECUTE PROCEDURE notify_event_store();


-- Create table to track the current state of projections.
CREATE TABLE projection_state
(
    projection_name          TEXT PRIMARY KEY,

    -- Track the information of the last event.
    last_processed_event_id  TEXT,
    last_processed_timestamp TIMESTAMP WITH TIME ZONE,
    aggregate_version        BIGINT                   NOT NULL DEFAULT 0,

    -- Track the global position
    global_position          NUMERIC                  NOT NULL DEFAULT 0,

    updated_at               TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);


-- Create table to store unique constraints.
CREATE TABLE unique_constraint
(
    field              CITEXT NOT NULL,
    value              CITEXT NOT NULL,
    owner_aggregate_id TEXT,
    created_at         TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (field, value)
);

-- Create index to delete by owner and field.
CREATE INDEX idx_unique_constraint_owner_aggregate_id_field ON unique_constraint (owner_aggregate_id, field);

-- Create keys table to store keys for encryption.
CREATE TABLE keys
(
    owner_id   TEXT PRIMARY KEY,
    key        bytea NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


-- Create table to store event journal lookup.
-- Think of it as a transactional consistent projection for certain properties of the event journal.
CREATE TABLE event_journal_lookup
(
    id                   TEXT PRIMARY KEY,
    owner_aggregate_id   TEXT NOT NULL,
    owner_aggregate_type TEXT NOT NULL,
    field_name           TEXT NOT NULL,
    field_value          TEXT NOT NULL,
    created_at           TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_event_journal_lookup_aggregate_type_field_name_field_value
    ON event_journal_lookup (owner_aggregate_type, field_name, field_value);

-- At most one field of a name per aggregate.
CREATE UNIQUE INDEX idx_event_journal_lookup_aggregate_id_field_name_unique
    ON event_journal_lookup (owner_aggregate_id, field_name);
