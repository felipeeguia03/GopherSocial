CREATE TABLE IF NOT EXISTS followers(
    user_id bigint not null,
    follower_id bigint not null,
    created_at timestamp(0) with time zone not null DEFAULT NOW(),

    primary key (user_id, follower_id),
    constraint fk_user foreign key (user_id) references users(id) ON DELETE CASCADE,
    constraint fk_follower foreign key (follower_id) references users(id) ON DELETE CASCADE
)