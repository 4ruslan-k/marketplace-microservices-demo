package applicationservices_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/pquerna/otp/totp"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"

	"authentication_service/config"

	mongoRepositories "authentication_service/internal/infrastructure/repositories/mongodb"
	fixtures "authentication_service/internal/test/fixtures"
	storage "authentication_service/pkg/storage/mongo"
	testUtils "authentication_service/pkg/testutils"

	"authentication_service/internal/application/dto"
	applicationServices "authentication_service/internal/application/services"
	socialAccountEntity "authentication_service/internal/domain/entities/social_account"
	userEntity "authentication_service/internal/domain/entities/user"
	domainServices "authentication_service/internal/domain/services"

	repo "authentication_service/internal/ports/repositories"
	middlewares "authentication_service/internal/transport/http/middlewares"
	mocks "authentication_service/mocks/pkg/messaging/nats"
)

var mongoURI string

func TestMain(m *testing.M) {
	_, uri, err := testUtils.InitializeMongoContainer(context.Background())
	mongoURI = uri
	if err != nil {
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func NewTestConfigWithDockerizedMongo(t *testing.T) *config.Config {
	t.Helper()
	testConf, err := config.NewTestConfig()
	if err != nil {
		t.Fatal(err)
	}
	testConf.MongoURI = mongoURI
	testConf.MongoDatabaseName = "test"
	if err != nil {
		t.Fatal(err)
	}
	return testConf
}

func NewTestSessionMiddleware(conf *config.Config, mongo *mongo.Database) *middlewares.Session {
	sessionStore := middlewares.NewSessionStore(mongo, conf)
	return middlewares.NewSession(sessionStore)
}

func NewTestApplicationService(
	conf *config.Config,
	mongo *mongo.Database,
	logger zerolog.Logger,
	t *testing.T,
) (applicationServices.UserApplicationService, repo.UserRepository, repo.AuthenticationRepository) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockNatsClient := mocks.NewMockNatsClient(mockCtrl)
	mockNatsClient.EXPECT().CreateStream(gomock.Any(), gomock.Any()).MinTimes(0)
	mockNatsClient.EXPECT().PublishMessage(gomock.Any(), gomock.Any()).MinTimes(0)
	userRepository := mongoRepositories.NewUserRepository(mongo, logger)
	authenticationRepository := mongoRepositories.NewAuthenticationRepository(mongo, logger)
	authenticationDomainService := domainServices.NewAuthenticationService(logger, authenticationRepository)
	userDomainService := domainServices.NewUserService(logger, authenticationDomainService, userRepository)
	applicationService := applicationServices.NewUserApplicationService(
		userRepository,
		authenticationRepository,
		logger,
		userDomainService,
		authenticationDomainService,
		mockNatsClient,
	)
	return applicationService, userRepository, authenticationRepository
}

func TestUserApplicationService_CreateUser(t *testing.T) {
	testConf := NewTestConfigWithDockerizedMongo(t)
	logger := zerolog.New(os.Stdout).Level(zerolog.Disabled)
	mongo := storage.NewMongoClient(logger, testConf)
	applicationService, _, _ := NewTestApplicationService(testConf, mongo, logger, t)

	email := fixtures.GenerateRandomEmail()
	cases := []struct {
		name   string
		input  userEntity.CreateUserParams
		output *dto.UserOutput
		expErr error
	}{
		{
			name:  "valid_input",
			input: userEntity.CreateUserParams{Name: "Joe", Password: "joe", Email: email},
			output: &dto.UserOutput{
				Name:  "Joe",
				Email: email,
			},
			expErr: nil,
		},
		{
			name:   "duplicated_email",
			input:  userEntity.CreateUserParams{Name: "Joe", Password: "joe", Email: email},
			output: nil,
			expErr: domainServices.ErrorEmailIsTaken,
		},
		{
			name:   "invalid_input_wrong_email_format",
			input:  userEntity.CreateUserParams{Name: "Joe", Password: "joe", Email: "example@"},
			output: nil,
			expErr: userEntity.ErrInvalidEmailFormat,
		},
		{
			name:   "invalid_input_empty_password",
			input:  userEntity.CreateUserParams{Name: "Joe", Password: "", Email: "example@gmail.com"},
			output: &dto.UserOutput{},
			expErr: domainServices.ErrInvalidPassword,
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			u, err := applicationService.CreateUser(context.Background(), tCase.input)
			if tCase.expErr != nil {
				require.ErrorContains(t, err, tCase.expErr.Error())
				return
			}
			require.NoError(t, err)
			require.NotNil(t, u)
			require.Equal(t, tCase.output.Email, u.Email)
			require.Equal(t, tCase.output.Name, u.Name)

		})
	}
}

func TestUserApplicationService_UpdateUser(t *testing.T) {
	testConf := NewTestConfigWithDockerizedMongo(t)
	logger := zerolog.New(os.Stdout).Level(zerolog.Disabled)
	mongo := storage.NewMongoClient(logger, testConf)

	applicationService, userRepository, _ := NewTestApplicationService(testConf, mongo, logger, t)

	type caseType struct {
		name             string
		generateTestData func()
		input            dto.UpdateUserInput
		output           *dto.UserOutput
		expErr           error
	}

	cases := []func() caseType{
		func() caseType {
			userID := fixtures.GenerateUUID()
			return caseType{
				name:   "error_user_not_found",
				input:  dto.UpdateUserInput{ID: userID, Name: "John"},
				expErr: applicationServices.ErrUserNotFound,
			}
		},
		func() caseType {
			userID := fixtures.GenerateUUID()
			return caseType{
				name:   "ok_user_updated",
				input:  dto.UpdateUserInput{ID: userID, Name: "John"},
				expErr: nil,
				generateTestData: func() {
					fixtures.IngestUser(t, fixtures.CreateTestUser{ID: userID}, userRepository.Create)
				},
				output: &dto.UserOutput{ID: userID},
			}
		},
	}

	for _, tCase := range cases {
		testData := tCase()
		t.Run(testData.name, func(t *testing.T) {
			if testData.generateTestData != nil {
				testData.generateTestData()
			}
			u, err := applicationService.UpdateUser(context.Background(), testData.input)
			if testData.expErr != nil {
				require.ErrorContains(t, err, testData.expErr.Error())
				return
			}
			require.NoError(t, err)
			require.NotNil(t, u)
			require.Equal(t, testData.output.ID, u.ID)
		})
	}
}

func TestUserApplicationService_DeleteUser(t *testing.T) {
	testConf := NewTestConfigWithDockerizedMongo(t)
	logger := zerolog.New(os.Stdout).Level(zerolog.Disabled)
	mongo := storage.NewMongoClient(logger, testConf)

	applicationService, userRepository, _ := NewTestApplicationService(testConf, mongo, logger, t)

	type caseType struct {
		name             string
		generateTestData func()
		input            dto.DeleteUserInput
		expErr           error
	}

	cases := []func() caseType{
		func() caseType {
			userID := fixtures.GenerateUUID()
			return caseType{
				name:   "error_user_not_found",
				input:  dto.DeleteUserInput{ID: userID},
				expErr: applicationServices.ErrUserNotFound,
			}
		},
		func() caseType {
			userID := fixtures.GenerateUUID()
			return caseType{
				name:   "ok_user_deleted",
				input:  dto.DeleteUserInput{ID: userID},
				expErr: nil,
				generateTestData: func() {
					fixtures.IngestUser(t, fixtures.CreateTestUser{ID: userID}, userRepository.Create)
				},
			}
		},
	}

	for _, tCase := range cases {
		testData := tCase()
		t.Run(testData.name, func(t *testing.T) {
			if testData.generateTestData != nil {
				testData.generateTestData()
			}
			err := applicationService.DeleteUser(context.Background(), testData.input)
			if testData.expErr != nil {
				require.ErrorContains(t, err, testData.expErr.Error())
				return
			}
			require.NoError(t, err)
			user, err := userRepository.GetByID(context.Background(), testData.input.ID)
			require.Nil(t, user)
			require.NoError(t, err)
		})
	}
}

func TestUserApplicationService_LoginWithEmailAndPassword(t *testing.T) {
	testConf := NewTestConfigWithDockerizedMongo(t)
	logger := zerolog.New(os.Stdout).Level(zerolog.Disabled)
	mongo := storage.NewMongoClient(logger, testConf)
	applicationService, userRepository, authenticationRepository := NewTestApplicationService(testConf, mongo, logger, t)

	type in struct {
		email        string
		password     string
		isMfaEnabled bool
	}

	type caseType struct {
		seedUser fixtures.CreateTestUser
		name     string
		input    in
		output   *dto.LoginOutput
		expErr   error
	}

	cases := []func() caseType{
		func() caseType {
			randomUser := fixtures.GenerateUserEntity(t, fixtures.CreateTestUser{})
			return caseType{
				name: "err_user_not_found",
				input: in{
					email:    randomUser.Email(),
					password: "wgwgwe#32r2",
				},
				output:   &dto.LoginOutput{Email: randomUser.Email()},
				seedUser: fixtures.CreateTestUser{},
				expErr:   applicationServices.ErrNoUserByEmail,
			}
		},
		func() caseType {
			email := fixtures.GenerateRandomEmail()
			return caseType{
				name: "err_invalid_password",
				input: in{
					email:    email,
					password: "wgwgwe#32r2",
				},
				seedUser: fixtures.CreateTestUser{
					Name:     "s",
					Email:    email,
					Password: "ss",
				},
				expErr: applicationServices.ErrInvalidCredentials,
			}
		},
		func() caseType {
			password := fixtures.GenerateRandomPassword()
			user := fixtures.GenerateUserEntity(t, fixtures.CreateTestUser{Password: password})
			name := "Sam"
			return caseType{
				name: "valid_login",
				input: in{
					email:    user.Email(),
					password: password,
				},
				output: &dto.LoginOutput{Email: user.Email(), Name: name},
				seedUser: fixtures.CreateTestUser{
					Name:     name,
					Email:    user.Email(),
					Password: password,
				},
			}
		},
		func() caseType {
			password := fixtures.GenerateRandomPassword()
			user := fixtures.GenerateUserEntity(t, fixtures.CreateTestUser{Password: password})
			name := "Sam"
			return caseType{
				name: "valid_login_mfa_enabled",
				input: in{
					email:        user.Email(),
					password:     password,
					isMfaEnabled: true,
				},
				output: &dto.LoginOutput{Email: "", Name: ""},
				seedUser: fixtures.CreateTestUser{
					Name:         name,
					Email:        user.Email(),
					Password:     password,
					IsMfaEnabled: true,
				},
			}
		},
	}

	for _, tCase := range cases {
		tCase := tCase()
		t.Run(tCase.name, func(t *testing.T) {
			fixtures.IngestUser(t, tCase.seedUser, userRepository.Create)
			u, err := applicationService.LoginWithEmailAndPassword(
				context.Background(),
				tCase.input.email,
				tCase.input.password,
			)
			if tCase.expErr != nil {
				require.ErrorContains(t, err, tCase.expErr.Error())
				return
			}
			require.NoError(t, err)
			require.NotNil(t, u)
			require.Equal(t, tCase.output.Email, u.Email)
			require.Equal(t, tCase.output.Name, u.Name)

			if tCase.input.isMfaEnabled {
				token, err := authenticationRepository.GetPasswordVerificationTokenByID(context.Background(), u.PasswordVerificationTokenID)
				require.NoError(t, err)
				require.Equal(t, false, token.HasExpired(time.Now()))
			}
		})
	}
}

func TestUserApplicationService_LoginWithTotpCode(t *testing.T) {
	testConf := NewTestConfigWithDockerizedMongo(t)
	logger := zerolog.New(os.Stdout)
	mongo := storage.NewMongoClient(logger, testConf)

	applicationService, userRepository, authenticationRepository := NewTestApplicationService(testConf, mongo, logger, t)

	type in struct {
		userID                      string
		otpCode                     string
		passwordVerificationTokenID string
	}

	type caseType struct {
		name             string
		generateTestData func()
		expErr           error
		in               in
	}

	cases := []func() caseType{
		func() caseType {
			userID := fixtures.GenerateUUID()
			return caseType{
				name:   "err_user_not_found",
				in:     in{userID: userID},
				expErr: applicationServices.ErrTotpCodeNotValid,
			}
		},
		func() caseType {
			userID := fixtures.GenerateUUID()
			return caseType{
				name:   "error_not_valid_code_empty_secret",
				in:     in{userID: userID},
				expErr: applicationServices.ErrTotpCodeNotValid,
				generateTestData: func() {
					fixtures.IngestUser(t, fixtures.CreateTestUser{ID: userID}, userRepository.Create)
				},
			}
		},
		func() caseType {
			userID := fixtures.GenerateUUID()
			key, _ := totp.Generate(totp.GenerateOpts{
				Issuer:      "test",
				AccountName: "test",
			})
			secret := key.Secret()
			return caseType{
				name:   "error_not_valid_code_valid_secret",
				in:     in{userID: userID, otpCode: "634212"},
				expErr: applicationServices.ErrTotpCodeNotValid,
				generateTestData: func() {
					fixtures.IngestUser(t, fixtures.CreateTestUser{ID: userID, TotpSecret: secret, IsMfaEnabled: true}, userRepository.Create)
				},
			}
		},
		func() caseType {
			userID := fixtures.GenerateUUID()
			passwordVerificationTokenID := fixtures.GenerateUUID()
			key, _ := totp.Generate(totp.GenerateOpts{
				Issuer:      "test",
				AccountName: "test",
			})
			secret := key.Secret()
			code, err := totp.GenerateCode(secret, time.Now())
			if err != nil {
				t.Fatal(err)
			}
			return caseType{
				name:   "ok_login_with_totp",
				in:     in{userID: userID, otpCode: code, passwordVerificationTokenID: passwordVerificationTokenID},
				expErr: nil,
				generateTestData: func() {
					fixtures.IngestUser(t, fixtures.CreateTestUser{ID: userID, TotpSecret: secret}, userRepository.Create)
					fixtures.IngestPasswordVerificationToken(t, fixtures.CreatePasswordVerificationToken{
						ID:        passwordVerificationTokenID,
						UserID:    userID,
						CreatedAt: time.Now(),
						ExpiresAt: time.Now().Add(time.Minute * 5),
					},
						authenticationRepository.SavePasswordVerificationToken)
				},
			}
		},
		func() caseType {
			userID := fixtures.GenerateUUID()
			passwordVerificationTokenID := fixtures.GenerateUUID()
			key, _ := totp.Generate(totp.GenerateOpts{
				Issuer:      "test",
				AccountName: "test",
			})
			secret := key.Secret()
			return caseType{
				name:   "error_wrong_code",
				in:     in{userID: userID, otpCode: "535356", passwordVerificationTokenID: passwordVerificationTokenID},
				expErr: applicationServices.ErrTotpCodeNotValid,
				generateTestData: func() {
					fixtures.IngestUser(t, fixtures.CreateTestUser{ID: userID, TotpSecret: secret}, userRepository.Create)
					fixtures.IngestPasswordVerificationToken(t, fixtures.CreatePasswordVerificationToken{
						ID:        passwordVerificationTokenID,
						UserID:    userID,
						CreatedAt: time.Now(),
						ExpiresAt: time.Now().Add(time.Minute * 5),
					},
						authenticationRepository.SavePasswordVerificationToken)
				},
			}
		},
		func() caseType {
			userID := fixtures.GenerateUUID()
			passwordVerificationTokenID := fixtures.GenerateUUID()
			key, _ := totp.Generate(totp.GenerateOpts{
				Issuer:      "test",
				AccountName: "test",
			})
			secret := key.Secret()
			code, err := totp.GenerateCode(secret, time.Now())
			if err != nil {
				t.Fatal(err)
			}
			return caseType{
				name:   "error_expired_password_verification_token",
				in:     in{userID: userID, otpCode: code, passwordVerificationTokenID: passwordVerificationTokenID},
				expErr: applicationServices.ErrTotpCodeNotValid,
				generateTestData: func() {
					fixtures.IngestUser(t, fixtures.CreateTestUser{ID: userID, TotpSecret: secret}, userRepository.Create)
					fixtures.IngestPasswordVerificationToken(t, fixtures.CreatePasswordVerificationToken{
						ID:        passwordVerificationTokenID,
						UserID:    userID,
						CreatedAt: time.Now(),
						ExpiresAt: time.Now().Add(-time.Minute * 5),
					},
						authenticationRepository.SavePasswordVerificationToken)
				},
			}
		},
	}

	for _, tCase := range cases {
		testData := tCase()
		t.Run(testData.name, func(t *testing.T) {
			if testData.generateTestData != nil {
				testData.generateTestData()
			}
			loginOutput, err := applicationService.LoginWithTotpCode(
				context.Background(),
				testData.in.passwordVerificationTokenID,
				testData.in.otpCode,
			)
			if testData.expErr != nil {
				require.ErrorContains(t, err, testData.expErr.Error())
				return
			}
			require.NoError(t, err)
			require.Equal(t, testData.in.userID, loginOutput.ID)
		})
	}
}

func TestUserApplicationService_SocialLogin(t *testing.T) {
	testConf := NewTestConfigWithDockerizedMongo(t)
	logger := zerolog.New(os.Stdout)
	mongo := storage.NewMongoClient(logger, testConf)

	applicationService, userRepository, _ := NewTestApplicationService(testConf, mongo, logger, t)

	type caseType struct {
		seedUser fixtures.CreateTestUser
		name     string
		input    dto.SocialLoginInput
		output   *dto.UserOutput
		expErr   error
	}

	cases := []func() caseType{
		func() caseType {
			randomUser := fixtures.GenerateUserEntity(t, fixtures.CreateTestUser{})
			return caseType{
				name: "ok_user_created_and_logged_in",
				input: dto.SocialLoginInput{
					Provider:  "google",
					Email:     randomUser.Email(),
					Name:      randomUser.Name(),
					UserID:    "s",
					AvatarURL: "",
				},
				output: &dto.UserOutput{Email: randomUser.Email(), Name: randomUser.Name()},
				expErr: nil,
			}
		},
		func() caseType {
			email := fixtures.GenerateRandomEmail()
			name := fixtures.GenerateRandomName()

			return caseType{
				name: "ok_already_created_user_logged_in",
				input: dto.SocialLoginInput{
					Provider:  "google",
					Email:     email,
					Name:      name,
					UserID:    "s",
					AvatarURL: "",
				},
				output:   &dto.UserOutput{Email: email, Name: name},
				expErr:   nil,
				seedUser: fixtures.CreateTestUser{Email: email, Name: name},
			}
		},
		func() caseType {
			name := fixtures.GenerateRandomName()
			return caseType{
				name: "err_empty_email",
				input: dto.SocialLoginInput{
					Provider:  "google",
					Email:     "",
					Name:      name,
					UserID:    "s",
					AvatarURL: "",
				},
				output: nil,
				expErr: userEntity.ErrInvalidEmailFormat,
			}
		},
		func() caseType {
			name := fixtures.GenerateRandomName()
			return caseType{
				name: "err_empty_provider",
				input: dto.SocialLoginInput{
					Provider:  "",
					Email:     "example@gmail.com",
					Name:      name,
					UserID:    "s",
					AvatarURL: "",
				},
				output: nil,
				expErr: socialAccountEntity.ErrInvalidProvider,
			}
		},
	}

	for _, tCase := range cases {
		tCase := tCase()
		t.Run(tCase.name, func(t *testing.T) {
			fixtures.IngestUser(t, tCase.seedUser, userRepository.Create)
			u, err := applicationService.SocialLogin(
				context.Background(),
				tCase.input,
			)
			if tCase.expErr != nil {
				require.ErrorContains(t, err, tCase.expErr.Error())
				return
			}
			require.NoError(t, err)
			require.NotNil(t, u)
			require.Equal(t, tCase.output.Email, u.Email)
			require.Equal(t, tCase.output.Name, u.Name)
		})
	}
}

func TestUserApplicationService_GetUserByID(t *testing.T) {
	testConf := NewTestConfigWithDockerizedMongo(t)
	logger := zerolog.New(os.Stdout)
	mongo := storage.NewMongoClient(logger, testConf)

	applicationService, userRepository, _ := NewTestApplicationService(testConf, mongo, logger, t)

	type caseType struct {
		seedUser fixtures.CreateTestUser
		name     string
		input    string
		output   *dto.UserOutput
		expErr   error
	}

	cases := []func() caseType{
		func() caseType {
			return caseType{
				name:   "ok_user_not_found",
				input:  "123",
				output: nil,
				expErr: nil,
			}
		},
		func() caseType {
			seedUser := fixtures.GenerateUserEntity(t, fixtures.CreateTestUser{})
			return caseType{
				name:     "ok_user_found",
				input:    seedUser.ID(),
				output:   &dto.UserOutput{ID: seedUser.ID(), Email: seedUser.Email(), Name: seedUser.Name()},
				expErr:   nil,
				seedUser: fixtures.CreateTestUser{ID: seedUser.ID(), Email: seedUser.Email(), Name: seedUser.Name()},
			}
		},
	}

	for _, tCase := range cases {
		tCase := tCase()
		t.Run(tCase.name, func(t *testing.T) {
			fixtures.IngestUser(t, tCase.seedUser, userRepository.Create)
			u, err := applicationService.GetUserByID(
				context.Background(),
				tCase.input,
			)
			if tCase.expErr != nil {
				require.ErrorContains(t, err, tCase.expErr.Error())
				return
			}
			require.NoError(t, err)
			require.Equal(t, tCase.output, u)
		})
	}
}

func TestUserApplicationService_ChangeCurrentPassword(t *testing.T) {
	testConf := NewTestConfigWithDockerizedMongo(t)
	logger := zerolog.New(os.Stdout)
	mongo := storage.NewMongoClient(logger, testConf)

	applicationService, userRepository, _ := NewTestApplicationService(testConf, mongo, logger, t)

	type caseType struct {
		name             string
		generateTestData func()
		input            dto.ChangeCurrentPasswordInput
		expErr           error
	}

	cases := []func() caseType{
		func() caseType {
			userID := fixtures.GenerateUUID()
			return caseType{
				name: "err_passwords_not_match",
				input: dto.ChangeCurrentPasswordInput{
					UserID:                  userID,
					CurrentPassword:         "123",
					NewPassword:             "1234",
					NewPasswordConfirmation: "12345",
				},
				expErr: applicationServices.ErrPasswordsDoNotMatch,
			}
		},
		func() caseType {
			userID := fixtures.GenerateUUID()
			newPass := "1234"

			return caseType{
				name: "err_user_not_found",
				input: dto.ChangeCurrentPasswordInput{
					UserID:                  userID,
					CurrentPassword:         "123",
					NewPassword:             newPass,
					NewPasswordConfirmation: newPass,
				},
				expErr: applicationServices.ErrNoUserByID,
			}
		},
		func() caseType {
			userID := fixtures.GenerateUUID()
			password := "h5h54h$%H45h45h4h"
			newPass := "1234"
			return caseType{
				name: "err_invalid_current_password",
				input: dto.ChangeCurrentPasswordInput{
					UserID:                  userID,
					CurrentPassword:         "123",
					NewPassword:             newPass,
					NewPasswordConfirmation: newPass,
				},
				expErr: applicationServices.ErrInvalidCredentials,
				generateTestData: func() {
					fixtures.IngestUser(t, fixtures.CreateTestUser{ID: userID, Password: password}, userRepository.Create)
				},
			}
		},
		func() caseType {
			userID := fixtures.GenerateUUID()
			password := "h5h54h$%H45h45h4h"
			newPass := "1234"
			return caseType{
				name: "err_invalid_current_password",
				input: dto.ChangeCurrentPasswordInput{
					UserID:                  userID,
					CurrentPassword:         password,
					NewPassword:             newPass,
					NewPasswordConfirmation: newPass,
				},
				expErr: nil,
				generateTestData: func() {
					fixtures.IngestUser(t, fixtures.CreateTestUser{ID: userID, Password: password}, userRepository.Create)
				},
			}
		},
	}

	for _, tCase := range cases {
		testData := tCase()
		t.Run(testData.name, func(t *testing.T) {
			if testData.generateTestData != nil {
				testData.generateTestData()
			}
			err := applicationService.ChangeCurrentPassword(
				context.Background(),
				testData.input,
			)
			if testData.expErr != nil {
				require.ErrorContains(t, err, testData.expErr.Error())
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestUserApplicationService_GenerateTotpSetup(t *testing.T) {
	testConf := NewTestConfigWithDockerizedMongo(t)
	logger := zerolog.New(os.Stdout)
	mongo := storage.NewMongoClient(logger, testConf)

	applicationService, userRepository, _ := NewTestApplicationService(testConf, mongo, logger, t)

	type caseType struct {
		name             string
		generateTestData func()
		userID           string
		expErr           error
	}

	cases := []func() caseType{
		func() caseType {
			userID := fixtures.GenerateUUID()
			return caseType{
				name:   "err_user_not_found",
				userID: userID,
				expErr: applicationServices.ErrNoUserByID,
			}
		},
		func() caseType {
			userID := fixtures.GenerateUUID()
			return caseType{
				name:   "ok_generate_totp_setup",
				userID: userID,
				expErr: nil,
				generateTestData: func() {
					fixtures.IngestUser(t, fixtures.CreateTestUser{ID: userID}, userRepository.Create)
				},
			}
		},
	}

	for _, tCase := range cases {
		testData := tCase()
		t.Run(testData.name, func(t *testing.T) {
			if testData.generateTestData != nil {
				testData.generateTestData()
			}
			totpSetupInfo, err := applicationService.GenerateTotpSetup(
				context.Background(),
				testData.userID,
			)
			if testData.expErr != nil {
				require.ErrorContains(t, err, testData.expErr.Error())
				return
			}
			require.NoError(t, err)
			require.NotEmpty(t, totpSetupInfo.Image)
			require.NotEmpty(t, totpSetupInfo.Secret)
			user, err := userRepository.GetByID(context.Background(), testData.userID)
			require.NoError(t, err)
			require.Equal(t, totpSetupInfo.Secret, user.MfaSettings().TotpSecret())
			require.Equal(t, false, user.MfaSettings().IsMfaEnabled())
		})
	}
}

func TestUserApplicationService_EnableTotp(t *testing.T) {
	testConf := NewTestConfigWithDockerizedMongo(t)
	logger := zerolog.New(os.Stdout)
	mongo := storage.NewMongoClient(logger, testConf)

	applicationService, userRepository, _ := NewTestApplicationService(testConf, mongo, logger, t)

	type caseType struct {
		name             string
		generateTestData func()
		userID           string
		otpCode          string
		expErr           error
	}

	cases := []func() caseType{
		func() caseType {
			userID := fixtures.GenerateUUID()
			return caseType{
				name:   "err_user_not_found",
				userID: userID,
				expErr: applicationServices.ErrNoUserByID,
			}
		},
		func() caseType {
			userID := fixtures.GenerateUUID()
			return caseType{
				name:   "error_not_valid_code_empty_secret",
				userID: userID,
				expErr: applicationServices.ErrTotpCodeNotValid,
				generateTestData: func() {
					fixtures.IngestUser(t, fixtures.CreateTestUser{ID: userID}, userRepository.Create)
				},
			}
		},
		func() caseType {
			userID := fixtures.GenerateUUID()
			key, _ := totp.Generate(totp.GenerateOpts{
				Issuer:      "test",
				AccountName: "test",
			})
			secret := key.Secret()
			return caseType{
				name:    "error_not_valid_code_valid_secret",
				userID:  userID,
				otpCode: "634212",
				expErr:  applicationServices.ErrTotpCodeNotValid,
				generateTestData: func() {
					fixtures.IngestUser(t, fixtures.CreateTestUser{ID: userID, TotpSecret: secret}, userRepository.Create)
				},
			}
		},
		func() caseType {
			userID := fixtures.GenerateUUID()
			key, _ := totp.Generate(totp.GenerateOpts{
				Issuer:      "test",
				AccountName: "test",
			})
			secret := key.Secret()
			code, err := totp.GenerateCode(secret, time.Now())
			if err != nil {
				t.Fatal(err)
			}
			return caseType{
				name:    "ok_enable_totp",
				userID:  userID,
				otpCode: code,
				expErr:  nil,
				generateTestData: func() {
					fixtures.IngestUser(t, fixtures.CreateTestUser{ID: userID, TotpSecret: secret}, userRepository.Create)
				},
			}
		},
	}

	for _, tCase := range cases {
		testData := tCase()
		t.Run(testData.name, func(t *testing.T) {
			if testData.generateTestData != nil {
				testData.generateTestData()
			}
			err := applicationService.EnableTotp(
				context.Background(),
				testData.userID,
				testData.otpCode,
			)
			if testData.expErr != nil {
				require.ErrorContains(t, err, testData.expErr.Error())
				return
			}
			require.NoError(t, err)
			user, err := userRepository.GetByID(context.Background(), testData.userID)
			require.NoError(t, err)
			require.Equal(t, true, user.MfaSettings().IsMfaEnabled())
		})
	}
}

func TestUserApplicationService_DisableTotp(t *testing.T) {
	testConf := NewTestConfigWithDockerizedMongo(t)
	logger := zerolog.New(os.Stdout)
	mongo := storage.NewMongoClient(logger, testConf)

	applicationService, userRepository, _ := NewTestApplicationService(testConf, mongo, logger, t)

	type caseType struct {
		name             string
		generateTestData func()
		userID           string
		otpCode          string
		expErr           error
	}

	cases := []func() caseType{
		func() caseType {
			userID := fixtures.GenerateUUID()
			return caseType{
				name:   "err_user_not_found",
				userID: userID,
				expErr: applicationServices.ErrNoUserByID,
			}
		},
		func() caseType {
			userID := fixtures.GenerateUUID()
			return caseType{
				name:   "error_not_valid_code_empty_secret",
				userID: userID,
				expErr: applicationServices.ErrTotpCodeNotValid,
				generateTestData: func() {
					fixtures.IngestUser(t, fixtures.CreateTestUser{ID: userID, IsMfaEnabled: true}, userRepository.Create)
				},
			}
		},
		func() caseType {
			userID := fixtures.GenerateUUID()
			key, _ := totp.Generate(totp.GenerateOpts{
				Issuer:      "test",
				AccountName: "test",
			})
			secret := key.Secret()
			return caseType{
				name:    "error_not_valid_code_valid_secret",
				userID:  userID,
				otpCode: "634212",
				expErr:  applicationServices.ErrTotpCodeNotValid,
				generateTestData: func() {
					fixtures.IngestUser(t, fixtures.CreateTestUser{ID: userID, TotpSecret: secret, IsMfaEnabled: true}, userRepository.Create)
				},
			}
		},
		func() caseType {
			userID := fixtures.GenerateUUID()
			key, _ := totp.Generate(totp.GenerateOpts{
				Issuer:      "test",
				AccountName: "test",
			})
			secret := key.Secret()
			code, err := totp.GenerateCode(secret, time.Now())
			if err != nil {
				t.Fatal(err)
			}
			return caseType{
				name:    "ok_disable_totp",
				userID:  userID,
				otpCode: code,
				expErr:  nil,
				generateTestData: func() {
					fixtures.IngestUser(t, fixtures.CreateTestUser{ID: userID, TotpSecret: secret, IsMfaEnabled: true}, userRepository.Create)
				},
			}
		},
	}

	for _, tCase := range cases {
		testData := tCase()
		t.Run(testData.name, func(t *testing.T) {
			if testData.generateTestData != nil {
				testData.generateTestData()
			}
			err := applicationService.DisableTotp(
				context.Background(),
				testData.userID,
				testData.otpCode,
			)
			if testData.expErr != nil {
				require.ErrorContains(t, err, testData.expErr.Error())
				return
			}
			require.NoError(t, err)
			user, err := userRepository.GetByID(context.Background(), testData.userID)
			require.NoError(t, err)
			require.Equal(t, false, user.MfaSettings().IsMfaEnabled())
		})
	}
}
