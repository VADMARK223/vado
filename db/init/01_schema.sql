create table tasks
(
    id          serial primary key,
    name        varchar(255)                        not null,
    description text,
    created_at  timestamp default CURRENT_TIMESTAMP not null,
    completed   boolean   default false,
    updated_at  timestamp default CURRENT_TIMESTAMP not null
);

comment on table tasks is 'Task table';