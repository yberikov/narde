CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users
(
    id         UUID PRIMARY KEY                  DEFAULT uuid_generate_v4(),
    username   VARCHAR(255) UNIQUE      NOT NULL,
    email      VARCHAR(255) UNIQUE      NOT NULL,
    password   TEXT                     NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
