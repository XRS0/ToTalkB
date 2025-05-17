CREATE TABLE users (
    id            serial PRIMARY KEY,
    login         varchar(255) NOT NULL UNIQUE,
    password_hash varchar(255) NOT NULL,
    name          varchar(255) NOT NULL,
    role          varchar(50)  NOT NULL
);

CREATE TABLE chats (
    id          serial PRIMARY KEY,
    name        varchar(255) NOT NULL,
    created_by  integer NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE chat_members (
    user_id integer NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    chat_id integer NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, chat_id)
);

CREATE TABLE messages (
    id         serial PRIMARY KEY,
    chat_id    integer NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    sender_id  integer NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at timestamp NOT NULL DEFAULT now(),
    content    text NOT NULL
);