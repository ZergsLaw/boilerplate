package migration

import zergrepo "github.com/ZergsLaw/zerg-repo"

var (
	Migrations = []zergrepo.Migrate{
		UserTable,
		SessionTable,
		NotificationTable,
		RecoveryCodeTable,
	}
)
