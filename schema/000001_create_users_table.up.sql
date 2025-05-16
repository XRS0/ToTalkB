CREATE TABLE users
(
    id            serial       NOT NULL UNIQUE,
    login         varchar(255) NOT NULL UNIQUE,
    password_hash varchar(255) NOT NULL,
    name          varchar(255) NOT NULL,
    role          varchar(50)  NOT NULL
);