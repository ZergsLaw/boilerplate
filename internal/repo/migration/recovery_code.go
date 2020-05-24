package migration

import zergrepo "github.com/ZergsLaw/zerg-repo"

const (
	createRecoveryCodeT = `
create table recovery_code
(
    id         serial,
    email      text                    not null,
    code       text                    not null,
    created_at timestamp default now() not null,

    foreign key (email) references users(email) on delete cascade,
    unique (code),
    primary key (id)
);
`
	dropRecoveryCodeT = `
drop table recovery_code;
`
)

var RecoveryCodeTable = zergrepo.Migrate{
	Version: 4,
	Up:      zergrepo.Query(createRecoveryCodeT),
	Down:    zergrepo.Query(dropRecoveryCodeT),
}
