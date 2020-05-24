package app_test

//func TestApp_StartWALNotification(t *testing.T) {
//	t.Parallel()
//
//	application, mocks, shutdown := initTest(t)
//	defer shutdown()
//
//	msg := app.Message{
//		Kind:    app.Welcome,
//		Content: "Welcome",
//	}
//
//	mocks.userRepo.EXPECT().UserByID(gomock.Any(), taskNotification.UserID).Return(&user1, nil).Times(3)
//	mocks.wal.EXPECT().NotificationTask(gomock.Any()).Return(&taskNotification, nil).Times(3)
//	mocks.notification.EXPECT().Notification(user1.Email, msg).Return(nil).Times(2)
//	mocks.wal.EXPECT().DeleteTaskNotification(gomock.Any(), taskNotification.ID).Return(nil)
//	mocks.wal.EXPECT().DeleteTaskNotification(gomock.Any(), taskNotification.ID).Return(errAny)
//	mocks.notification.EXPECT().Notification(user1.Email, msg).Return(errAny)
//	mocks.wal.EXPECT().NotificationTask(gomock.Any()).Return(nil, app.ErrNotFound)
//	mocks.wal.EXPECT().NotificationTask(gomock.Any()).Return(nil, errAny)
//
//	testCases := []struct {
//		name string
//		want error
//	}{
//		{"err delete task", errAny},
//		{"err send notification", errAny},
//		{"err get notification task", errAny},
//	}
//
//	for _, tc := range testCases {
//		tc := tc
//		t.Run(tc.name, func(t *testing.T) {
//			err := application.StartWALNotification(ctx)
//			assert.Equal(t, tc.want, err)
//		})
//	}
//}
