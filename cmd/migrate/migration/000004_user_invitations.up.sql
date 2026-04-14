CREATE TABLE IF NOT EXISTS user_invitations (
    token bytea primary key,
    user_id bigint not null
)