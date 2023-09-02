package applicationservices_test

import (
	"context"
	"notification/internal/test/fixtures"
	"os"
	"testing"
	"time"

	testUtils "notification/pkg/testutils"

	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"

	"notification/config"
	notificationEntity "notification/internal/domain/entities/notification"
	notificationRepo "notification/internal/repositories/notification"
	notificationRepoPg "notification/internal/repositories/notification/pg"
	"notification/migrate/migrations"
	pgStorage "shared/storage/pg"

	applicationServices "notification/internal/services"
	mocks "notification/mocks/pkg/messaging/nats"
)

var pgDSN string

func TestMain(m *testing.M) {
	_, uri, err := testUtils.InitializePGContainer(context.Background())
	pgDSN = uri
	if err != nil {
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func NewTestConfigWithDockerizePG(t *testing.T) *config.Config {
	ctx := context.Background()
	t.Helper()
	testConf, err := config.NewTestConfig()
	if err != nil {
		t.Fatal(err)
	}
	testConf.PgSDN = pgDSN
	if err != nil {
		t.Fatal(err)
	}
	db := pgStorage.NewClientWithDSN(zerolog.Logger{}, pgDSN, false)

	migrator := migrate.NewMigrator(db, migrations.Migrations)
	err = migrator.Init(ctx)
	if err != nil {
		t.Fatal()
	}

	_, err = migrator.Migrate(ctx)
	if err != nil {
		t.Fatal()
	}

	return testConf
}

func NewTestApplicationService(
	conf *config.Config,
	pg *bun.DB,
	logger zerolog.Logger,
	t *testing.T,
) (applicationServices.NotificationApplicationService, notificationRepo.NotificationsRepository) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockNatsClient := mocks.NewMockNatsClient(mockCtrl)
	mockNatsClient.EXPECT().CreateStream(gomock.Any(), gomock.Any()).MinTimes(0)
	mockNatsClient.EXPECT().PublishMessageEphemeral(gomock.Any(), gomock.Any()).MinTimes(0)
	notificationRepository := notificationRepoPg.NewNotificationRepository(pg, logger)
	applicationService := applicationServices.NewNotificationApplicationService(
		notificationRepository,
		logger,
		mockNatsClient,
	)
	return applicationService, notificationRepository
}

func TestUserApplicationService_CreateUserNotification(t *testing.T) {
	t.Parallel()
	testConf := NewTestConfigWithDockerizePG(t)
	logger := zerolog.New(os.Stdout).Level(zerolog.Disabled)
	pg := pgStorage.NewClient(logger, pgStorage.Config{DSN: testConf.PgSDN})

	notificationApplicationService, notificationRepository := NewTestApplicationService(testConf, pg, logger, t)

	type caseType struct {
		name             string
		args             notificationEntity.CreateUserNotificationParams
		expErr           error
		generateTestData func()
	}

	testCases := []func() caseType{
		func() caseType {
			return caseType{
				name: "ok",
				args: notificationEntity.CreateUserNotificationParams{
					UserID:             fixtures.GenerateUUID(),
					NotificationTypeID: notificationEntity.MFADisabledNotification.TypeID(),
				},
			}
		},
		func() caseType {
			return caseType{
				name: "error_empty_user_id",
				args: notificationEntity.CreateUserNotificationParams{
					NotificationTypeID: notificationEntity.MFADisabledNotification.TypeID(),
				},
				expErr: notificationEntity.ErrUserIDEmpty,
			}
		},
	}

	for _, tCase := range testCases {
		tCase := tCase()
		t.Run(tCase.name, func(t *testing.T) {
			t.Parallel()
			if tCase.generateTestData != nil {
				tCase.generateTestData()
			}
			err := notificationApplicationService.CreateUserNotification(context.Background(), tCase.args)
			if tCase.expErr != nil {
				require.ErrorContains(t, err, tCase.expErr.Error())
				return
			}
			require.NoError(t, err)
			notificationsList, err := notificationRepository.GetByUserID(context.Background(), tCase.args.UserID)
			require.NoError(t, err)
			require.Equal(t, 1, len(notificationsList))
		})
	}
}

func TestUserApplicationService_GetNotificationsByUserID(t *testing.T) {
	t.Parallel()
	testConf := NewTestConfigWithDockerizePG(t)
	logger := zerolog.New(os.Stdout).Level(zerolog.Disabled)
	pg := pgStorage.NewClient(logger, pgStorage.Config{DSN: testConf.PgSDN})

	notificationApplicationService, notificationRepository := NewTestApplicationService(testConf, pg, logger, t)

	type args struct {
		userID           string
		notificationType string
		viewedAt         time.Time
	}

	type caseType struct {
		name             string
		userID           string
		args             args
		expErr           error
		generateTestData func()
	}

	testCases := []func() caseType{
		func() caseType {
			return caseType{
				name:   "error_empty_type_id",
				userID: "",
				expErr: applicationServices.ErrInvalidUserID,
			}
		},
		func() caseType {
			userID := fixtures.GenerateUUID()
			viewedAt := time.Now()
			notificationType := notificationEntity.MFADisabledNotification.TypeID()
			return caseType{
				name: "create_notification",
				args: args{
					userID:           userID,
					notificationType: notificationType,
					viewedAt:         viewedAt,
				},
				userID: userID,
				expErr: nil,
				generateTestData: func() {
					fixtures.InsertUserNotification(
						t,
						fixtures.CreateTestUserNotification{
							NotificationTypeID: notificationType,
							UserID:             userID,
							CreatedAt:          time.Now(),
							ViewedAt:           viewedAt,
						},
						notificationRepository.CreateUserNotification,
					)
				},
			}
		},
	}

	for _, tCase := range testCases {
		tCase := tCase()
		t.Run(tCase.name, func(t *testing.T) {
			t.Parallel()
			if tCase.generateTestData != nil {
				tCase.generateTestData()
			}
			notifications, err := notificationApplicationService.GetNotificationsByUserID(context.Background(), tCase.userID)
			if tCase.expErr != nil {
				require.ErrorContains(t, err, tCase.expErr.Error())
				return
			}
			require.NoError(t, err)
			require.Equal(t, 1, len(notifications.Notifications))

			notification := notifications.Notifications[0]
			require.Equal(t, tCase.args.notificationType, notification.NotificationTypeID)
			require.Equal(t, tCase.args.viewedAt.Format(time.RFC3339), notification.ViewedAt.Local().Format(time.RFC3339))
		})
	}
}

func TestUserApplicationService_ViewNotification(t *testing.T) {
	t.Parallel()
	testConf := NewTestConfigWithDockerizePG(t)
	logger := zerolog.New(os.Stdout).Level(zerolog.Disabled)
	pg := pgStorage.NewClient(logger, pgStorage.Config{DSN: testConf.PgSDN})

	notificationApplicationService, notificationRepository := NewTestApplicationService(testConf, pg, logger, t)

	type args struct {
		userID             string
		userNotificationID string
	}

	type caseType struct {
		name             string
		args             args
		expErr           error
		generateTestData func()
	}

	testCases := []func() caseType{
		func() caseType {
			userID := fixtures.GenerateUUID()
			userNotificationID := fixtures.GenerateUUID()
			return caseType{
				name:   "error_notification_not_found",
				args:   args{userID: userID, userNotificationID: userNotificationID},
				expErr: applicationServices.ErrNotificationNotFound,
			}
		},
		func() caseType {
			userID := fixtures.GenerateUUID()
			userNotificationID := fixtures.GenerateUUID()
			return caseType{
				name:   "view_one_notification",
				args:   args{userID: userID, userNotificationID: userNotificationID},
				expErr: nil,
				generateTestData: func() {
					fixtures.InsertUserNotification(
						t,
						fixtures.CreateTestUserNotification{
							ID:                 userNotificationID,
							NotificationTypeID: notificationEntity.MFAEnabledNotification.TypeID(),
							UserID:             userID,
							CreatedAt:          time.Now(),
						},
						notificationRepository.CreateUserNotification,
					)
				},
			}
		},
	}

	for _, tCase := range testCases {
		tCase := tCase()
		t.Run(tCase.name, func(t *testing.T) {
			t.Parallel()
			if tCase.generateTestData != nil {
				tCase.generateTestData()
			}
			err := notificationApplicationService.ViewNotification(context.Background(), tCase.args.userID, tCase.args.userNotificationID)
			if tCase.expErr != nil {
				require.ErrorContains(t, err, tCase.expErr.Error())
				return
			}
			require.NoError(t, err)

			notification, err := notificationRepository.GetByUserIDAndUserNotificationID(context.Background(), tCase.args.userID, tCase.args.userNotificationID)
			require.NoError(t, err)
			require.Equal(t, notification.IsZero(), false)
			require.Equal(t, notification.UpdatedAt().IsZero(), false)
			require.Equal(t, notification.ViewedAt().IsZero(), false)
		})
	}
}

func TestUserApplicationService_DeleteUserNotification(t *testing.T) {
	t.Parallel()
	testConf := NewTestConfigWithDockerizePG(t)
	logger := zerolog.New(os.Stdout).Level(zerolog.Disabled)
	pg := pgStorage.NewClient(logger, pgStorage.Config{DSN: testConf.PgSDN})

	notificationApplicationService, notificationRepository := NewTestApplicationService(testConf, pg, logger, t)

	type args struct {
		userID             string
		userNotificationID string
	}

	type caseType struct {
		name             string
		args             args
		expErr           error
		generateTestData func()
	}

	testCases := []func() caseType{
		func() caseType {
			userID := fixtures.GenerateUUID()
			userNotificationID := fixtures.GenerateUUID()
			return caseType{
				name:   "notification_not_found",
				args:   args{userID: userID, userNotificationID: userNotificationID},
				expErr: nil,
			}
		},
		func() caseType {
			userID := fixtures.GenerateUUID()
			userNotificationID := fixtures.GenerateUUID()
			return caseType{
				name:   "delete_one_notification",
				args:   args{userID: userID, userNotificationID: userNotificationID},
				expErr: nil,
				generateTestData: func() {
					fixtures.InsertUserNotification(
						t,
						fixtures.CreateTestUserNotification{
							ID:                 userNotificationID,
							NotificationTypeID: notificationEntity.MFADisabledNotification.TypeID(),
							UserID:             userID,
							CreatedAt:          time.Now(),
						},
						notificationRepository.CreateUserNotification,
					)
				},
			}
		},
	}

	for _, tCase := range testCases {
		tCase := tCase()
		t.Run(tCase.name, func(t *testing.T) {
			t.Parallel()
			if tCase.generateTestData != nil {
				tCase.generateTestData()
			}
			err := notificationApplicationService.DeleteUserNotification(context.Background(), tCase.args.userID, tCase.args.userNotificationID)
			if tCase.expErr != nil {
				require.ErrorContains(t, err, tCase.expErr.Error())
				return
			}
			require.NoError(t, err)

			notification, err := notificationRepository.GetByUserIDAndUserNotificationID(context.Background(), tCase.args.userID, tCase.args.userNotificationID)
			require.NoError(t, err)
			require.Equal(t, notification.IsZero(), true)
		})
	}
}

func TestUserApplicationService_ViewAllNotifications(t *testing.T) {
	testConf := NewTestConfigWithDockerizePG(t)
	logger := zerolog.New(os.Stdout).Level(zerolog.Disabled)
	pg := pgStorage.NewClient(logger, pgStorage.Config{DSN: testConf.PgSDN})

	notificationApplicationService, notificationRepository := NewTestApplicationService(testConf, pg, logger, t)

	type args struct {
		userID string
	}

	type caseType struct {
		name             string
		args             args
		expErr           error
		generateTestData func()
	}

	testCases := []func() caseType{
		func() caseType {
			userID := fixtures.GenerateUUID()
			userNotificationID := fixtures.GenerateUUID()
			return caseType{
				name:   "view_all_notifications",
				args:   args{userID: userID},
				expErr: nil,
				generateTestData: func() {
					fixtures.InsertUserNotification(
						t,
						fixtures.CreateTestUserNotification{
							ID:                 userNotificationID,
							NotificationTypeID: notificationEntity.MFADisabledNotification.TypeID(),
							UserID:             userID,
							CreatedAt:          time.Now(),
						},
						notificationRepository.CreateUserNotification,
					)

				},
			}
		},
	}

	for _, tCase := range testCases {
		tCase := tCase()
		t.Run(tCase.name, func(t *testing.T) {
			t.Parallel()
			if tCase.generateTestData != nil {
				tCase.generateTestData()
			}
			err := notificationApplicationService.ViewAllNotifications(context.Background(), tCase.args.userID)
			if tCase.expErr != nil {
				require.ErrorContains(t, err, tCase.expErr.Error())
				return
			}
			require.NoError(t, err)

			notifications, err := notificationRepository.GetByUserID(context.Background(), tCase.args.userID)
			require.NoError(t, err)

			for _, notification := range notifications {
				require.Equal(t, notification.IsZero(), false)
				require.Equal(t, notification.ViewedAt().IsZero(), false)
				require.Equal(t, notification.UpdatedAt().IsZero(), false)

			}
		})
	}
}
