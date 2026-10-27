CREATE TABLE roles (
    id SMALLINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

INSERT INTO roles (name)
VALUES
    ('customer'),
    ('admin'),
    ('superadmin');

CREATE TABLE users (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL,
    email CITEXT NOT NULL UNIQUE,
    password_hash BYTEA NOT NULL,
    activated BOOLEAN NOT NULL DEFAULT FALSE,

    role_id SMALLINT NOT NULL
        REFERENCES roles(id),

    created_at TIMESTAMPTZ(0) NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ(0) NOT NULL DEFAULT NOW(),

    version INTEGER NOT NULL DEFAULT 1
        CHECK (version > 0)
);
