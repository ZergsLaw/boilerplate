package migration

import zergrepo "github.com/ZergsLaw/zerg-repo"

const (
	createUserT = `
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
);`
	dropUserT = `
drop table users;
`
)

var UserTable = zergrepo.Migrate{
	Version: 1,
	Up:      zergrepo.Query(createUserT),
	Down:    zergrepo.Query(dropUserT),
}
