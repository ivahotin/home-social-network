create table if not exists followers (
    follower_id     bigint not null,
    user_id         bigint not null,
    is_active       boolean not null,
    created_at      timestamp default current_timestamp(),
    primary key (user_id, follower_id)
) engine = InnoDB;