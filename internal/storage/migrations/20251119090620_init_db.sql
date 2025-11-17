-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE questions (
                           id SERIAL PRIMARY KEY,
                           text VARCHAR NOT NULL,
                           created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE answers (
                         id SERIAL PRIMARY KEY,
                         question_id INTEGER NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
                         user_id UUID NOT NULL DEFAULT uuid_generate_v4(),
                         text VARCHAR NOT NULL,
                         created_at TIMESTAMP NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS answers;
DROP TABLE IF EXISTS questions;
