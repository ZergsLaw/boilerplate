--up
create table recovery_code
(
    id         serial,
    email      text                    not null,
    code       text                    not null,
    created_at timestamp default now() not null,

    foreign key (email) references users (email) on delete cascade,
    unique (code),
    primary key (id)
);


--down
drop table recovery_code;
