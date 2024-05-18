-- Define roles table
CREATE TABLE roles (
    role_id BIGSERIAL PRIMARY KEY,
    role_name VARCHAR NOT NULL UNIQUE
);

INSERT INTO roles (role_name) VALUES ('admin'), ('user'), ('moderator');

-- Define users table with bigint serial primary key
CREATE TABLE users (
    user_id BIGSERIAL PRIMARY KEY,
    username VARCHAR UNIQUE NOT NULL,
    hashed_password VARCHAR NOT NULL,
    fullname VARCHAR NOT NULL,
    email VARCHAR UNIQUE NOT NULL,
    is_email_verified BOOLEAN DEFAULT false,
    password_changed_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    role_id BIGINT REFERENCES roles(role_id) ON DELETE SET NULL
);

-- Define accounts table to track user balances
CREATE TABLE accounts (
    id BIGSERIAL PRIMARY KEY,
    owner_user_id BIGINT REFERENCES users(user_id) ON DELETE CASCADE,
    balance BIGINT NOT NULL,
    unit VARCHAR NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (owner_user_id, unit)
);

-- Define entries table to track account transactions
CREATE TABLE entries (
    id BIGSERIAL PRIMARY KEY,
    account_id BIGINT REFERENCES accounts(id) ON DELETE CASCADE,
    amount BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Define transfers table to track fund transfers between accounts
CREATE TABLE transfers (
    id BIGSERIAL PRIMARY KEY,
    from_account_id BIGINT REFERENCES accounts(id) ON DELETE CASCADE,
    to_account_id BIGINT REFERENCES accounts(id) ON DELETE CASCADE,
    amount BIGINT NOT NULL CHECK (amount > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Define sessions table to manage user sessions
CREATE TABLE sessions (
    id UUID PRIMARY KEY,
    user_id BIGINT REFERENCES users(user_id) ON DELETE CASCADE,
    refresh_token VARCHAR NOT NULL,
    user_agent VARCHAR NOT NULL,
    client_ip VARCHAR NOT NULL,
    is_blocked BOOLEAN NOT NULL DEFAULT false,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Define verify_emails table to manage email verification requests
CREATE TABLE verify_emails (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(user_id) ON DELETE CASCADE,
    email VARCHAR NOT NULL,
    secret_code VARCHAR NOT NULL,
    is_used BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expired_at TIMESTAMPTZ NOT NULL DEFAULT (CURRENT_TIMESTAMP + INTERVAL '15 minutes')
);

-- Define locations table to store geographical locations
CREATE TABLE locations (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    notes TEXT, 
    coords GEOGRAPHY(POINT, 4326) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Define resources table to track resources associated with locations
CREATE TABLE resources (
    id BIGSERIAL PRIMARY KEY,
    location_id BIGINT REFERENCES locations(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    notes TEXT, 
    owner_user_id BIGINT REFERENCES users(user_id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create index on frequently queried columns
CREATE INDEX ON accounts (owner_user_id);
CREATE INDEX ON entries (account_id);
CREATE INDEX ON transfers (from_account_id);
CREATE INDEX ON transfers (to_account_id);
CREATE INDEX ON locations USING GIST (coords);
