--up
create table sessions
(
    id         serial,
    user_id    integer   not null,
    token_id   text      not null,
    ip         inet      not null,
    user_agent text      not null default '',
    created_at timestamp not null default now(),
    is_logout  bool      not null default false,

    foreign key (user_id) references users on delete cascade,
    unique (token_id),
    primary key (id)
);


--down
drop table sessions;
