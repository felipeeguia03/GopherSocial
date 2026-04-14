CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users(
    id bigserial primary key,
    username varchar(255) not null,
    email citext unique not null,
    password bytea not null,
    created_at timestamp(0) with time zone not null DEFAULT NOW()
)