CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users
(
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    given_name TEXT NOT NULL,
    family_name TEXT NOT NULL,
    hashed_password TEXT NOT NULL,
    password_last_updated TIMESTAMP NOT NULL DEFAULT NOW(),
    email_is_confirmed BOOLEAN DEFAULT FALSE,
    email_confirmation_code uuid DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (username, email)
);

CREATE TABLE IF NOT EXISTS password_reset_tokens
(
    id uuid PRIMARY KEY NOT NULL,
    token uuid UNIQUE DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expiry TIMESTAMP NOT NULL DEFAULT NOW() + INTERVAL '1 HOUR'
);

