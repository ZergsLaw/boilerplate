--up
create table users
(
    id         serial,
    email      text      not null,
    username   text      not null,
    pass_hash  BYTEA,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),

    unique (username),
    unique (email),
    primary key (id)
);

--down
drop table users;