package mock

//go:generate mockgen -source=../app/app.go -aux_files github.com/zergslaw/boilerplate/internal/app=../app/user.go -destination mock.app.contracts.go -package mock
//go:generate mockgen -source=../app/user.go -destination=mock.user.contracts.go -package mock
//go:generate mockgen -source=../app/notification.go -destination=mock.notification.contracts.go -package mock
//go:generate mockgen -source=../app/wal.go -destination=mock.wal.contracts.go -package mock
