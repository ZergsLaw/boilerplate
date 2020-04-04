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
    user_id    integer                 not null,
    token_id   text                    not null,
    ip         inet                    not null,
    user_agent text                    not null default '',
    created_at timestamp default now() not null,
    is_logout  bool      default false not null,

    foreign key (user_id) references users on delete cascade,
    unique (token_id),
    primary key (id)
);

create table notifications
(
    id         serial,
    user_id    integer                 not null,
    kind       text                    not null,
    is_done    bool      default false not null,
    created_at timestamp default now() not null,
    exec_time  timestamp,

    foreign key (user_id) references users on delete cascade,
    primary key (id)
);

create table recovery_code
(
    id         serial,
    user_id    integer                 not null,
    code       text                    not null,
    created_at timestamp default now() not null,

    foreign key (user_id) references users on delete cascade,
    unique (code),
    primary key (id)
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
drop table sessions;
drop table notifications;
drop table recovery_code;
drop table users;
