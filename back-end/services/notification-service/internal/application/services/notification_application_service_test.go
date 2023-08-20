package applicationservices_test

import (
	"context"
	"notification_service/internal/test/fixtures"
	"os"
	"testing"
	"time"

	testUtils "notification_service/pkg/testutils"

	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"

	"notification_service/config"
	notificationEntity "notification_service/internal/domain/entities/notification"
	pgrepositories "notification_service/internal/infrastructure/repositories/pg/notification"
	repo "notification_service/internal/ports/repositories"
	"notification_service/migrate/migrations"
	pgStorage "notification_service/pkg/storage/pg"

	applicationServices "notification_service/internal/application/services"
	mocks "notification_service/mocks/pkg/messaging/nats"
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
) (applicationServices.NotificationApplicationService, repo.NotificationsRepository) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockNatsClient := mocks.NewMockNatsClient(mockCtrl)
	mockNatsClient.EXPECT().CreateStream(gomock.Any(), gomock.Any()).MinTimes(0)
	mockNatsClient.EXPECT().PublishMessageEphemeral(gomock.Any(), gomock.Any()).MinTimes(0)
	notificationRepository := pgrepositories.NewNotificationRepository(pg, logger)
	applicationService := applicationServices.NewNotificationApplicationService(
		notificationRepository,
		logger,
		mockNatsClient,
	)
	return applicationService, notificationRepository
}

func TestUserApplicationService_CreateUserNotification(t *testing.T) {
	testConf := NewTestConfigWithDockerizePG(t)
	logger := zerolog.New(os.Stdout).Level(zerolog.Disabled)
	pg := pgStorage.NewClient(logger, testConf)

	notificationApplicationService, notificationRepository := NewTestApplicationService(testConf, pg, logger, t)

	type caseType struct {
		name             string
		input            notificationEntity.CreateUserNotificationParams
		expErr           error
		generateTestData func()
		cleanup          func()
	}

	cases := []func() caseType{
		func() caseType {
			return caseType{
				name: "ok",
				input: notificationEntity.CreateUserNotificationParams{
					UserID:             fixtures.GenerateUUID(),
					NotificationTypeID: notificationEntity.MFADisabledNotification.TypeID(),
				},
				generateTestData: func() {
				},
				cleanup: func() {
				},
			}
		},
		func() caseType {
			return caseType{
				name: "error_empty_user_id",
				input: notificationEntity.CreateUserNotificationParams{
					NotificationTypeID: notificationEntity.MFADisabledNotification.TypeID(),
				},
				expErr: notificationEntity.ErrUserIDEmpty,
			}
		},
	}

	for _, tCase := range cases {
		tCase := tCase()
		t.Run(tCase.name, func(t *testing.T) {
			if tCase.generateTestData != nil {
				tCase.generateTestData()
			}
			err := notificationApplicationService.CreateUserNotification(context.Background(), tCase.input)
			if tCase.expErr != nil {
				require.ErrorContains(t, err, tCase.expErr.Error())
				return
			}
			require.NoError(t, err)
			notificationsList, err := notificationRepository.GetByUserID(context.Background(), tCase.input.UserID)
			require.NoError(t, err)
			require.Equal(t, 1, len(notificationsList))
			if tCase.cleanup != nil {
				tCase.cleanup()
			}
		})
	}
}

func TestUserApplicationService_GetNotificationsByUserID(t *testing.T) {
	testConf := NewTestConfigWithDockerizePG(t)
	logger := zerolog.New(os.Stdout).Level(zerolog.Disabled)
	pg := pgStorage.NewClient(logger, testConf)

	notificationApplicationService, notificationRepository := NewTestApplicationService(testConf, pg, logger, t)

	type Input struct {
		userID           string
		notificationType string
		viewedAt         time.Time
	}

	type caseType struct {
		name             string
		userID           string
		input            Input
		expErr           error
		generateTestData func()
		cleanup          func()
	}

	cases := []func() caseType{
		func() caseType {
			return caseType{
				name:   "err_empty_type_id",
				userID: "",
				expErr: applicationServices.ErrInvalidUserID,
			}
		},
		func() caseType {
			userID := fixtures.GenerateUUID()
			viewedAt := time.Now()
			notificationType := notificationEntity.MFADisabledNotification.TypeID()
			return caseType{
				name: "ok_create_notification",
				input: Input{
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
				cleanup: func() {

				},
			}
		},
	}

	for _, tCase := range cases {
		tCase := tCase()
		t.Run(tCase.name, func(t *testing.T) {
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
			require.Equal(t, tCase.input.notificationType, notification.NotificationTypeID)
			require.Equal(t, tCase.input.viewedAt.Format(time.RFC3339), notification.ViewedAt.Local().Format(time.RFC3339))

			if tCase.cleanup != nil {
				tCase.cleanup()
			}
		})
	}
}

func TestUserApplicationService_ViewNotification(t *testing.T) {
	testConf := NewTestConfigWithDockerizePG(t)
	logger := zerolog.New(os.Stdout).Level(zerolog.Disabled)
	pg := pgStorage.NewClient(logger, testConf)

	notificationApplicationService, notificationRepository := NewTestApplicationService(testConf, pg, logger, t)

	type Input struct {
		userID             string
		userNotificationID string
	}

	type caseType struct {
		name             string
		input            Input
		expErr           error
		generateTestData func()
		cleanup          func()
	}

	cases := []func() caseType{
		func() caseType {
			userID := fixtures.GenerateUUID()
			userNotificationID := fixtures.GenerateUUID()
			return caseType{
				name:   "err_notification_not_found",
				input:  Input{userID: userID, userNotificationID: userNotificationID},
				expErr: applicationServices.ErrNotificationNotFound,
			}
		},
		func() caseType {
			userID := fixtures.GenerateUUID()
			userNotificationID := fixtures.GenerateUUID()
			return caseType{
				name:   "ok_view_one_notification",
				input:  Input{userID: userID, userNotificationID: userNotificationID},
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
				cleanup: func() {
				},
			}
		},
	}

	for _, tCase := range cases {
		tCase := tCase()
		t.Run(tCase.name, func(t *testing.T) {
			if tCase.generateTestData != nil {
				tCase.generateTestData()
			}
			err := notificationApplicationService.ViewNotification(context.Background(), tCase.input.userID, tCase.input.userNotificationID)
			if tCase.expErr != nil {
				require.ErrorContains(t, err, tCase.expErr.Error())
				return
			}
			require.NoError(t, err)

			notification, err := notificationRepository.GetByUserIDAndUserNotificationID(context.Background(), tCase.input.userID, tCase.input.userNotificationID)
			require.NoError(t, err)
			require.Equal(t, notification.IsZero(), false)
			require.Equal(t, notification.UpdatedAt().IsZero(), false)
			require.Equal(t, notification.ViewedAt().IsZero(), false)

			if tCase.cleanup != nil {
				tCase.cleanup()
			}
		})
	}
}

func TestUserApplicationService_DeleteUserNotification(t *testing.T) {
	testConf := NewTestConfigWithDockerizePG(t)
	logger := zerolog.New(os.Stdout).Level(zerolog.Disabled)
	pg := pgStorage.NewClient(logger, testConf)

	notificationApplicationService, notificationRepository := NewTestApplicationService(testConf, pg, logger, t)

	type Input struct {
		userID             string
		userNotificationID string
	}

	type caseType struct {
		name             string
		input            Input
		expErr           error
		generateTestData func()
		cleanup          func()
	}

	cases := []func() caseType{
		func() caseType {
			userID := fixtures.GenerateUUID()
			userNotificationID := fixtures.GenerateUUID()
			return caseType{
				name:   "ok_notification_not_found",
				input:  Input{userID: userID, userNotificationID: userNotificationID},
				expErr: nil,
			}
		},
		func() caseType {
			userID := fixtures.GenerateUUID()
			userNotificationID := fixtures.GenerateUUID()
			return caseType{
				name:   "ok_delete_one_notification",
				input:  Input{userID: userID, userNotificationID: userNotificationID},
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
				cleanup: func() {
				},
			}
		},
	}

	for _, tCase := range cases {
		tCase := tCase()
		t.Run(tCase.name, func(t *testing.T) {
			if tCase.generateTestData != nil {
				tCase.generateTestData()
			}
			err := notificationApplicationService.DeleteUserNotification(context.Background(), tCase.input.userID, tCase.input.userNotificationID)
			if tCase.expErr != nil {
				require.ErrorContains(t, err, tCase.expErr.Error())
				return
			}
			require.NoError(t, err)

			notification, err := notificationRepository.GetByUserIDAndUserNotificationID(context.Background(), tCase.input.userID, tCase.input.userNotificationID)
			require.NoError(t, err)
			require.Equal(t, notification.IsZero(), true)

			if tCase.cleanup != nil {
				tCase.cleanup()
			}
		})
	}
}

func TestUserApplicationService_ViewAllNotifications(t *testing.T) {
	testConf := NewTestConfigWithDockerizePG(t)
	logger := zerolog.New(os.Stdout).Level(zerolog.Disabled)
	pg := pgStorage.NewClient(logger, testConf)

	notificationApplicationService, notificationRepository := NewTestApplicationService(testConf, pg, logger, t)

	type Input struct {
		userID string
	}

	type caseType struct {
		name             string
		input            Input
		expErr           error
		generateTestData func()
		cleanup          func()
	}

	cases := []func() caseType{
		func() caseType {
			userID := fixtures.GenerateUUID()
			userNotificationID := fixtures.GenerateUUID()
			return caseType{
				name:   "ok_view_all_notifications",
				input:  Input{userID: userID},
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
				cleanup: func() {
				},
			}
		},
	}

	for _, tCase := range cases {
		tCase := tCase()
		t.Run(tCase.name, func(t *testing.T) {
			if tCase.generateTestData != nil {
				tCase.generateTestData()
			}
			err := notificationApplicationService.ViewAllNotifications(context.Background(), tCase.input.userID)
			if tCase.expErr != nil {
				require.ErrorContains(t, err, tCase.expErr.Error())
				return
			}
			require.NoError(t, err)

			notifications, err := notificationRepository.GetByUserID(context.Background(), tCase.input.userID)
			require.NoError(t, err)

			for _, notification := range notifications {
				require.Equal(t, notification.IsZero(), false)
				require.Equal(t, notification.ViewedAt().IsZero(), false)
				require.Equal(t, notification.UpdatedAt().IsZero(), false)

			}

			if tCase.cleanup != nil {
				tCase.cleanup()
			}
		})
	}
}
