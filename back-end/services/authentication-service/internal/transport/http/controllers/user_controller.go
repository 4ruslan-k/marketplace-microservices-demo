package controllers

import (
	"authentication_service/config"
	applicationServices "authentication_service/internal/application/services"
	userEntity "authentication_service/internal/domain/entities/user"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/mongo/mongodriver"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"

	domainDto "authentication_service/internal/application/dto"
	httpDto "authentication_service/internal/transport/http/dto"
	httpErrors "authentication_service/pkg/errors/http"
)

type UserControllers struct {
	ApplicationService applicationServices.UserApplicationService
	Logger             zerolog.Logger
	Config             *config.Config
}

type UserOutput struct {
	User *domainDto.UserOutput `json:"user"`
}

func NewUserControllers(
	appService applicationServices.UserApplicationService,
	logger zerolog.Logger,
	config *config.Config,
) *UserControllers {
	return &UserControllers{
		ApplicationService: appService,
		Logger:             logger,
		Config:             config,
	}
}

func SetupSocialLogin(
	mb *mongo.Database,
	config *config.Config,
) {

	goth.UseProviders(
		github.New(config.SocialSignIn.GithubKey, config.SocialSignIn.GithubSecret, config.GatewayURL+"/v1/auth/social/github/callback"),
		google.New(config.SocialSignIn.GoogleKey, config.SocialSignIn.GoogleSecret, config.GatewayURL+"/v1/auth/social/google/callback"),
	)

	sessionsCollection := mb.Collection("sessions")
	store := mongodriver.NewStore(sessionsCollection, 3600, true, []byte(config.SessionSecret))

	gothic.Store = store

	authMap := make(map[string]string)
	authMap["github"] = "Github"
	authMap["google"] = "Google"
}

func (r *UserControllers) SocialLoginCallback(c *gin.Context) {
	gothicUser, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	r.Logger.Info().Interface("gothicUser", gothicUser).Msg("gothicUser")
	socialAccount := domainDto.SocialLoginInput{
		Provider:  gothicUser.Provider,
		Email:     gothicUser.Email,
		Name:      gothicUser.Name,
		UserID:    gothicUser.UserID,
		AvatarURL: gothicUser.AvatarURL,
	}

	userLoginOutput, err := r.ApplicationService.SocialLogin(c.Request.Context(), socialAccount)
	if err != nil {
		httpErrors.RespondWithError(c, err)
		return
	}

	if !userLoginOutput.IsMfaEnabled {
		session := sessions.Default(c)
		err = saveSession(session, userLoginOutput.ID)
		if err != nil {
			httpErrors.RespondWithError(c, err)
			return
		}
	}

	link := fmt.Sprintf("%s", r.Config.FrontendURL)
	if userLoginOutput.IsMfaEnabled {
		link = fmt.Sprintf("%s/signin", link)
		queryParams := url.Values{}
		queryParams.Add("isMfaFlow", "true")
		queryParams.Add("token", userLoginOutput.PasswordVerificationTokenID)
		link = link + "?" + queryParams.Encode()
	}

	c.Redirect(http.StatusPermanentRedirect, link)
}

func (r *UserControllers) SocialLogin(c *gin.Context) {
	// try to get the user without re-authenticating
	if gothicUser, err := gothic.CompleteUserAuth(c.Writer, c.Request); err == nil {
		socialAccount := domainDto.SocialLoginInput{
			Provider:  gothicUser.Provider,
			Email:     gothicUser.Email,
			Name:      gothicUser.Name,
			UserID:    gothicUser.UserID,
			AvatarURL: gothicUser.AvatarURL,
		}

		userLoginOutput, err := r.ApplicationService.SocialLogin(c.Request.Context(), socialAccount)
		if err != nil {
			httpErrors.RespondWithError(c, err)
			return
		}
		if !userLoginOutput.IsMfaEnabled {
			session := sessions.Default(c)
			err = saveSession(session, userLoginOutput.ID)
			if err != nil {
				httpErrors.RespondWithError(c, err)
				return
			}
		}

		link := fmt.Sprintf("%s", r.Config.FrontendURL)
		if userLoginOutput.IsMfaEnabled {
			link = fmt.Sprintf("%s/signin", link)
			queryParams := url.Values{}
			queryParams.Add("isMfaFlow", "true")
			queryParams.Add("token", userLoginOutput.PasswordVerificationTokenID)
			link = link + "?" + queryParams.Encode()
		}
		c.Redirect(http.StatusPermanentRedirect, link)

	} else {
		context := context.WithValue(c.Request.Context(), "provider", c.Params.ByName("provider"))
		gothic.BeginAuthHandler(c.Writer, c.Request.WithContext(context))
	}
	handleOkResponse(c)
}

func getUserIDFromSession(session sessions.Session) *string {
	sessionUserID := session.Get("user_id")
	if sessionUserID == nil {
		return nil
	}
	userID := sessionUserID.(string)
	return &userID
}

func (r *UserControllers) GetCurrentUser(c *gin.Context) {
	session := sessions.Default(c)
	userID := getUserIDFromSession(session)
	if userID == nil {
		handleResponseWithBody(c, httpDto.UserOutput{User: nil})
		return
	}
	userOutput, err := r.ApplicationService.GetUserByID(c.Request.Context(), *userID)
	if userOutput == nil {
		httpErrors.Unauthorized(c, "Not Authorized, no user found")
		return
	}
	if err != nil {
		httpErrors.RespondWithError(c, err)
		return
	}
	r.Logger.Info().Str("userId", userOutput.ID).Str("SessionID", session.ID()).Msg("UserControllers -> GetCurrentUser")
	handleResponseWithBody(c, httpDto.UserOutput{User: userOutput})
}

func (r *UserControllers) GetCurrentUserInternal(c *gin.Context) {
	session := sessions.Default(c)
	userID := getUserIDFromSession(session)
	sessionID := session.ID()

	if userID == nil {
		handleResponseWithBody(c, httpDto.UserWithSessionIDOutput{User: nil, SessionID: sessionID})
		return
	}
	userOutput, err := r.ApplicationService.GetUserByID(c.Request.Context(), *userID)
	if err != nil {
		httpErrors.RespondWithError(c, err)
		return
	}
	r.Logger.Info().
		Str("userId", userOutput.ID).
		Str("sessionID", sessionID).
		Msg("UserControllers -> GetCurrentUserInternal")

	handleResponseWithBody(c, httpDto.UserWithSessionIDOutput{User: userOutput, SessionID: sessionID})
}

func (r *UserControllers) GetUserByID(c *gin.Context) {
	userIDfromParams, found := c.Params.Get("userID")
	if found == false {
		httpErrors.RespondWithError(c, errors.New("user_id parameter not found"))
		return
	}
	session := sessions.Default(c)
	sessionUserID := session.Get("user_id").(string)
	if userIDfromParams != sessionUserID {
		httpErrors.Unauthorized(c, "You are not authorized to fetch this user")
		return
	}
	user, err := r.ApplicationService.GetUserByID(c.Request.Context(), sessionUserID)
	if err != nil {
		httpErrors.RespondWithError(c, err)
		return
	}
	handleResponseWithBody(c, httpDto.UserOutput{User: user})
}

// Creates a user from user input
func (h *UserControllers) CreateUser(c *gin.Context) {
	var createUserInput httpDto.CreateUserInput
	if err := c.ShouldBindJSON(&createUserInput); err != nil {
		httpErrors.BadRequest(c, err.Error())
		return
	}

	userOutput, err := h.ApplicationService.CreateUser(
		c.Request.Context(),
		userEntity.CreateUserParams{
			Name:          createUserInput.Name,
			Email:         createUserInput.Email,
			Password:      createUserInput.Password,
			SocialAccount: nil,
		},
	)
	if err != nil {
		httpErrors.RespondWithError(c, err)
		return
	}
	session := sessions.Default(c)
	err = saveSession(session, userOutput.ID)
	handleSuccessResponse(c, http.StatusCreated, "created")
}

// Updates a user
func (h *UserControllers) UpdateUser(c *gin.Context) {
	userIDfromParams := c.Param("userID")
	var updateUserInput httpDto.UpdateUserInput
	if err := c.ShouldBindJSON(&updateUserInput); err != nil {
		httpErrors.BadRequest(c, err.Error())
		return
	}
	session := sessions.Default(c)
	sessionUserID := session.Get("user_id").(string)
	if userIDfromParams != sessionUserID {
		httpErrors.Unauthorized(c, "You are not authorized to update this user")
		return
	}
	userID := session.Get("user_id").(string)
	_, err := h.ApplicationService.UpdateUser(
		c.Request.Context(),
		domainDto.UpdateUserInput{
			Name: updateUserInput.Name,
			ID:   userID,
		},
	)
	if err != nil {
		httpErrors.RespondWithError(c, err)
		return
	}

	handleOkResponse(c)
}

// Deletes a user
func (h *UserControllers) DeleteUser(c *gin.Context) {
	userIDfromParams := c.Param("userID")
	session := sessions.Default(c)
	sessionUserID := session.Get("user_id").(string)
	if userIDfromParams != sessionUserID {
		httpErrors.Unauthorized(c, "You are not authorized to delete this user")
		return
	}
	userID := session.Get("user_id").(string)
	err := h.ApplicationService.DeleteUser(
		c.Request.Context(),
		domainDto.DeleteUserInput{
			ID: userID,
		},
	)
	if err != nil {
		httpErrors.RespondWithError(c, err)
		return
	}

	handleOkResponse(c)
}

func saveSession(session sessions.Session, userID string) error {
	session.Set("user_id", userID)
	err := session.Save()
	if err != nil {
		return err
	}
	return nil
}

type LoginInput struct {
	Email    string `json:"username"  binding:"required"`
	Password string `json:"password"  binding:"required"`
}

// Logins with email and password
func (h *UserControllers) LoginWithEmailAndPassword(c *gin.Context) {
	var loginInput LoginInput
	if err := c.ShouldBindJSON(&loginInput); err != nil {
		httpErrors.BadRequest(c, err.Error())
		return
	}

	loginOutput, err := h.ApplicationService.LoginWithEmailAndPassword(
		c.Request.Context(),
		loginInput.Email,
		loginInput.Password,
	)

	if err != nil {
		httpErrors.RespondWithError(c, err)
		return
	}
	if !loginOutput.IsMfaEnabled {
		session := sessions.Default(c)
		err = saveSession(session, loginOutput.ID)
		if err != nil {
			httpErrors.RespondWithError(c, err)
			return
		}
	}

	handleResponseWithBody(c, httpDto.LoginOutput{LoginOutput: loginOutput})
}

type LoginWithTotpInput struct {
	PasswordVerificationTokenID string `json:"tokenId"  binding:"required"`
	Code                        string `json:"code"  binding:"required"`
}

// Logins with TOTP code, requires a password verification token
func (h *UserControllers) LoginWithTotpCode(c *gin.Context) {
	var loginInput LoginWithTotpInput
	if err := c.ShouldBindJSON(&loginInput); err != nil {
		httpErrors.BadRequest(c, err.Error())
		return
	}

	userOutput, err := h.ApplicationService.LoginWithTotpCode(
		c.Request.Context(),
		loginInput.PasswordVerificationTokenID,
		loginInput.Code,
	)

	if err != nil {
		httpErrors.RespondWithError(c, err)
		return
	}

	session := sessions.Default(c)
	err = saveSession(session, userOutput.ID)
	if err != nil {
		httpErrors.RespondWithError(c, err)
		return
	}
	handleResponseWithBody(c, userOutput)
}

func (h *UserControllers) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	err := session.Save()
	if err != nil {
		httpErrors.RespondWithError(c, err)
		return
	}
	handleOkResponse(c)
}

// Change Password
type ChangePasswordInput struct {
	CurrentPassword         string `json:"currentPassword"  binding:"required"`
	NewPassword             string `json:"newPassword"  binding:"required"`
	NewPasswordConfirmation string `json:"newPasswordConfirmation"  binding:"required"`
}

func (h *UserControllers) ChangeCurrentPassword(c *gin.Context) {
	var changePasswordInput ChangePasswordInput
	if err := c.ShouldBindJSON(&changePasswordInput); err != nil {
		httpErrors.BadRequest(c, err.Error())
		return
	}

	session := sessions.Default(c)
	sessionUserID := session.Get("user_id").(string)

	err := h.ApplicationService.ChangeCurrentPassword(
		c.Request.Context(),
		domainDto.ChangeCurrentPasswordInput{
			UserID:                  sessionUserID,
			CurrentPassword:         changePasswordInput.CurrentPassword,
			NewPassword:             changePasswordInput.NewPassword,
			NewPasswordConfirmation: changePasswordInput.NewPasswordConfirmation,
		},
	)

	if err != nil {
		httpErrors.RespondWithError(c, err)
		return
	}
	handleOkResponse(c)
}

// Generates MFA TOTP setup - user tries to enable TOTP based MFA
func (h *UserControllers) GenerateTotpSetup(c *gin.Context) {
	session := sessions.Default(c)
	sessionUserID := session.Get("user_id").(string)

	totpSetup, err := h.ApplicationService.GenerateTotpSetup(
		c.Request.Context(),
		sessionUserID,
	)

	if err != nil {
		httpErrors.RespondWithError(c, err)
		return
	}
	handleResponseWithBody(c, httpDto.GenerateTotpSetupOutput{TotpSetup: httpDto.TotpSetup{
		Image: totpSetup.Image,
	}})
}

type EnableTotpMfaInput struct {
	Otp string `json:"code" binding:"required"`
}

// Enables TOTP MFA by validating the OTP code
func (h *UserControllers) EnableTotpMfa(c *gin.Context) {
	var enableTotpMfaInput EnableTotpMfaInput
	if err := c.ShouldBindJSON(&enableTotpMfaInput); err != nil {
		httpErrors.BadRequest(c, err.Error())
		return
	}
	session := sessions.Default(c)
	sessionUserID := session.Get("user_id").(string)

	err := h.ApplicationService.EnableTotp(
		c.Request.Context(),
		sessionUserID,
		enableTotpMfaInput.Otp,
	)

	if err != nil {
		httpErrors.RespondWithError(c, err)
		return
	}
	handleOkResponse(c)
}

// Disables TOTP MFA by validating the OTP code
func (h *UserControllers) DisableTotpMfa(c *gin.Context) {
	var enableTotpMfaInput EnableTotpMfaInput
	if err := c.ShouldBindJSON(&enableTotpMfaInput); err != nil {
		httpErrors.BadRequest(c, err.Error())
		return
	}
	session := sessions.Default(c)
	sessionUserID := session.Get("user_id").(string)

	err := h.ApplicationService.DisableTotp(
		c.Request.Context(),
		sessionUserID,
		enableTotpMfaInput.Otp,
	)

	if err != nil {
		httpErrors.RespondWithError(c, err)
		return
	}
	handleOkResponse(c)
}
