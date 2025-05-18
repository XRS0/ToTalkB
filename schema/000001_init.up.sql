CREATE TABLE IF NOT EXISTS users (
    id            serial PRIMARY KEY,  -- integer
    login         varchar(255) NOT NULL UNIQUE,
    password_hash varchar(255) NOT NULL,
    name          varchar(255) NOT NULL,
    role          varchar(50)  NOT NULL
);

CREATE TABLE IF NOT EXISTS chats (
    id          serial PRIMARY KEY,  -- integer
    name        varchar(255) NOT NULL,
    created_by  integer NOT NULL REFERENCES users(id) ON DELETE CASCADE  -- integer
);

CREATE TABLE IF NOT EXISTS chat_members (
    user_id integer NOT NULL REFERENCES users(id) ON DELETE CASCADE,  -- integer
    chat_id integer NOT NULL REFERENCES chats(id) ON DELETE CASCADE,   -- integer
    PRIMARY KEY (user_id, chat_id)
);

CREATE TABLE IF NOT EXISTS messages (
    id         serial PRIMARY KEY,  -- integer
    chat_id    integer NOT NULL REFERENCES chats(id) ON DELETE CASCADE,  -- integer
    sender_id  integer NOT NULL REFERENCES users(id) ON DELETE CASCADE,  -- integer
    created_at timestamp NOT NULL DEFAULT now(),
    content    text NOT NULL
);

CREATE TABLE IF NOT EXISTS events (
    id UUID PRIMARY KEY,
    type VARCHAR(255) NOT NULL,
    source VARCHAR(255) NOT NULL,
    payload BYTEA NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
); 

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

CREATE TABLE IF NOT EXISTS notifications (
    id VARCHAR(36) PRIMARY KEY,
    user_id INTEGER NOT NULL,
    type VARCHAR(50) NOT NULL,
    payload JSONB NOT NULL,
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    scheduled_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_status ON notifications(status);
CREATE INDEX IF NOT EXISTS idx_notifications_scheduled_at ON notifications(scheduled_at); 