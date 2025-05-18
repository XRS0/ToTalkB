CREATE TABLE IF NOT EXISTS event_queues (
    id UUID PRIMARY KEY,
    event_id UUID NOT NULL REFERENCES events(id),
    user_id UUID NOT NULL,
    status VARCHAR(50) NOT NULL,
    position INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    UNIQUE(event_id, user_id)
);

CREATE INDEX idx_event_queues_event_id ON event_queues(event_id);
CREATE INDEX idx_event_queues_user_id ON event_queues(user_id);
CREATE INDEX idx_event_queues_status ON event_queues(status); 