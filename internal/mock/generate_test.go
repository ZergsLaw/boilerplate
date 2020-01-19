package mock

//go:generate mockgen -source=../app/app.go -destination=mock.go -package mock -mock_names App=App,Repo=Repo,Notification=Notification,Password=Password,Auth=Auth,OAuth=OAuth
