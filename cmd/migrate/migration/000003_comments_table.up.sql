CREATE TABLE IF NOT EXISTS comments (
    id bigserial primary key,
    user_id bigint not null,
    post_id bigint not null,
    content text not null,
    created_at timestamp(0) with time zone not null DEFAULT NOW(),
    constraint fk_user foreign key (user_id) references users(id),
    constraint fk_post foreign key (post_id) references posts(id)
)