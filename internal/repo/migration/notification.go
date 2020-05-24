package migration

import zergrepo "github.com/ZergsLaw/zerg-repo"

const (
	createNotificationT = `
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
`
	dropNotificationT = `
drop table notifications;
`
)

var NotificationTable = zergrepo.Migrate{
	Version: 3,
	Up:      zergrepo.Query(createNotificationT),
	Down:    zergrepo.Query(dropNotificationT),
}
