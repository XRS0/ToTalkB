CREATE TABLE users (
    id            serial PRIMARY KEY,  -- integer
    login         varchar(255) NOT NULL UNIQUE,
    password_hash varchar(255) NOT NULL,
    name          varchar(255) NOT NULL,
    role          varchar(50)  NOT NULL
);

CREATE TABLE chats (
    id          serial PRIMARY KEY,  -- integer
    name        varchar(255) NOT NULL,
    created_by  integer NOT NULL REFERENCES users(id) ON DELETE CASCADE  -- integer
);

CREATE TABLE chat_members (
    user_id integer NOT NULL REFERENCES users(id) ON DELETE CASCADE,  -- integer
    chat_id integer NOT NULL REFERENCES chats(id) ON DELETE CASCADE,   -- integer
    PRIMARY KEY (user_id, chat_id)
);

CREATE TABLE messages (
    id         serial PRIMARY KEY,  -- integer
    chat_id    integer NOT NULL REFERENCES chats(id) ON DELETE CASCADE,  -- integer
    sender_id  integer NOT NULL REFERENCES users(id) ON DELETE CASCADE,  -- integer
    created_at timestamp NOT NULL DEFAULT now(),
    content    text NOT NULL
);