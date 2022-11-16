CREATE
EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE lifecycle AS ENUM ('created', 'past');

CREATE TABLE broadcasts
(
    id          UUID                     NOT NULL PRIMARY KEY,
    name        VARCHAR(255)             NOT NULL,
    owner       VARCHAR(100)             NOT NULL,
    description text                     NOT NULL,
    previewUrl  VARCHAR(255) DEFAULT '',
    streamKey   VARCHAR(150)             NOT NULL,
    life        lifecycle                NOT NULL,
    start_time  timestamp with time zone NOT NULL
);

CREATE TABLE images
(
    id   UUID REFERENCES broadcasts ON DELETE CASCADE,
    file bytea
);

CREATE TABLE files
(
    id           UUID PRIMARY KEY NOT NULL,
    broadcast_id UUID REFERENCES broadcasts (id) ON DELETE CASCADE,
    name         VARCHAR(255)     NOT NULL,
    type         VARCHAR(100) DEFAULT '',
    size         serial           NOT NULL,
    file         bytea            NOT NULL
);

CREATE TABLE messages
(
    id          UUID PRIMARY KEY         NOT NULL,
    channel     VARCHAR(36)              NOT NULL,
    username    VARCHAR(150)             NOT NULL,
    fullname    VARCHAR(200)             NOT NULL,
    avatar      VARCHAR(64)              NOT NULL,
    text        TEXT                     NOT NULL,
    time        timestamp with time zone NOT NULL,
    is_question BOOLEAN                  NOT NULL DEFAULT FALSE,
    is_anon     BOOLEAN                  NOT NULL DEFAULT FALSE,
    reactions   jsonb                    NOT NULL DEFAULT '{}'::jsonb
);

CREATE TABLE participants
(
    id       UUID PRIMARY KEY NOT NULL,
    channel  VARCHAR(36)      NOT NULL,
    username VARCHAR(200)     NOT NULL,
    fullname VARCHAR(200),
    email    VARCHAR(100)
);

CREATE TABLE users
(
    id       UUID         NOT NULL PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE
);

INSERT INTO users (id, username)
VALUES (uuid_generate_v4(), 'petrov'),
       (uuid_generate_v4(), 'ivanov'),
       (uuid_generate_v4(), 'sidotov'),
       (uuid_generate_v4(), 'pushkin');

CREATE TABLE live
(
    id          UUID         NOT NULL PRIMARY KEY UNIQUE,
    place       VARCHAR(100),
    description VARCHAR(255) NOT NULL,
    streamUrl   VARCHAR(255) NOT NULL
);

INSERT INTO live (id, place, description, streamUrl)
VALUES (uuid_generate_v4(), 'г. Москва', 'Смоленская-Сенная', 		'/camera/3stream'),
       (uuid_generate_v4(), 'г. Москва', 'Воронцовcкая, д. 5/2', '/camera/1stream'),
       (uuid_generate_v4(), 'г. Москва', 'Воронцовская, д. 1/3', '/camera/4stream'),
       (uuid_generate_v4(), 'г. Москва', 'Петровский', '/camera/2stream'),
       (uuid_generate_v4(), 'г. Москва', 'Дорожный', '/camera/5stream'),
       (uuid_generate_v4(), 'г. Новосибирск', 'Ленина 20', '/camera/7stream'),
       (uuid_generate_v4(), 'г. Екатеринбург', 'Вайнера 40', '/camera/6stream'),
       (uuid_generate_v4(), 'г. Екатеринбург', 'Блюхера 53а', '/camera/ekat2');

CREATE TABLE stream
(
    id          UUID         NOT NULL PRIMARY KEY,
    username    VARCHAR(150) NOT NULL,
    description VARCHAR(255) DEFAULT ''
);

CREATE TABLE zoom
(
    id              UUID                     NOT NULL PRIMARY KEY,
    start_time      timestamp with time zone NOT NULL,
    email           VARCHAR(255)             NOT NULL,
    topic           TEXT,
    recording_count smallint DEFAULT 0,
    json            json                     NOT NULL
);