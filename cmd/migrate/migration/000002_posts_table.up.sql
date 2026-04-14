CREATE TABLE IF NOT EXISTS posts (
    id bigserial primary key,
    title varchar(255) not null,
    content text not null,
    user_id bigint not null,
    tags text[] not null DEFAULT '{}',
    created_at timestamp(0) with time zone not null DEFAULT NOW(),
    updated_at timestamp(0) with time zone not null DEFAULT NOW(),
    version int not null DEFAULT 0,
 
    constraint fk_user foreign key (user_id) references users(id)
)