create table if not exists friends (
    id           serial primary key,
    user_id      int not null,
    friend_id    int not null,
    unique key   first_user_id_second_user_id (first_user_id, second_user_id)
);