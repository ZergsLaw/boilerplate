--up
create table notifications
(
    id         serial,
    email      text                    not null,
    kind       text                    not null,
    is_done    bool      default false not null,
    created_at timestamp default now() not null,
    exec_time  timestamp,

    foreign key (email) references users(email) on delete cascade on update cascade,
    primary key (id)
);


--down
drop table notifications;
