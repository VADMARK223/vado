CREATE TABLE IF NOT EXISTS tasks
(
    id          serial primary key,
    name        varchar(255) not null,
    description text,
    created_at  timestamp default CURRENT_TIMESTAMP,
    completed   boolean   default false
);