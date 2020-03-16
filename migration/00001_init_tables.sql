-- +goose Up
-- SQL in this section is executed when the migration is applied.
create table users
(
    id         serial,
    email      text                    not null,
    username   text                    not null,
    pass_hash  BYTEA,
    created_at timestamp default now() not null,
    updated_at timestamp default now() not null,

    unique (username),
    unique (email),
    primary key (id)
);

create table sessions
(
    id         serial,
    user_id    integer references users (id) on delete cascade,
    token_id   text                    not null,
    ip         inet                    not null,
    user_agent text                    not null default '',
    created_at timestamp default now() not null,
    is_logout  bool      default false not null,

    unique (token_id),
    primary key (id)
);

create table notifications
(
    id         serial,
    user_id    integer references users on delete cascade,
    kind       text                    not null,
    is_done    bool      default false not null,
    created_at timestamp default now() not null,
    exec_time  timestamp,

    primary key (id)
);

create table recovery_code
(
    id         serial,
    email      text references users (email) on delete cascade,
    code       text                    not null,
    created_at timestamp default now() not null,

    unique (email),
    unique (code),
    primary key (id)
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
drop table sessions;
drop table notifications;
drop table recovery_code;
drop table users;
