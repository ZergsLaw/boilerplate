package app_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/zergslaw/boilerplate/internal/app"
)

func TestApp_StartWAL(t *testing.T) {
	t.Parallel()

	application, _, _, _, mockWal, mockNotification, shutdown := initTest(t)
	defer shutdown()

	mockWal.EXPECT().NotificationTask(gomock.Any()).Return(&taskNotification, nil).Times(3)
	mockNotification.EXPECT().Notification(taskNotification.Email, app.Welcome).Return(nil).Times(2)
	mockWal.EXPECT().DeleteTaskNotification(gomock.Any(), taskNotification.ID).Return(nil)
	mockWal.EXPECT().DeleteTaskNotification(gomock.Any(), taskNotification.ID).Return(errAny)
	mockNotification.EXPECT().Notification(taskNotification.Email, app.Welcome).Return(errAny)
	mockWal.EXPECT().NotificationTask(gomock.Any()).Return(nil, app.ErrNotFound)
	mockWal.EXPECT().NotificationTask(gomock.Any()).Return(nil, errAny)

	testCases := []struct {
		name string
		want error
	}{
		{"err delete task", errAny},
		{"err send notification", errAny},
		{"err get notification task", errAny},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := application.StartWALNotification(ctx)
			assert.Equal(t, tc.want, err)
		})
	}
}
