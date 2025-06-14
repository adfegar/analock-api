package services

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/adfer-dev/analock-api/auth"
	"github.com/adfer-dev/analock-api/models"
)

// AuthService struct
type AuthService struct {
	googleValidator GoogleTokenValidator
	AppTokenManager auth.TokenManager
	userService     UserService
	tokenService    TokenService
	extLoginService ExternalLoginService
}

// AuthService constructor
func NewAuthService(
	googleValidator GoogleTokenValidator,
	appTokenManager auth.TokenManager,
	userService UserService,
	tokenService TokenService,
	extLoginService ExternalLoginService,
) *AuthService {
	return &AuthService{
		googleValidator: googleValidator,
		AppTokenManager: appTokenManager,
		userService:     userService,
		tokenService:    tokenService,
		extLoginService: extLoginService,
	}
}

// Request bodies
type UserAuthenticateBody struct {
	Email         string `json:"email" validate:"required,email"`
	UserName      string `json:"userName" validate:"required"`
	ProviderId    string `json:"providerId" validate:"required"`
	ProviderToken string `json:"providerToken" validate:"required,jwt"`
}

type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required,jwt"`
}

type RefreshTokenResponse struct {
	Token string `json:"token"`
}

// AuthService methods
func (s *AuthService) AuthenticateUser(authBody UserAuthenticateBody) (*models.Token, *models.Token, error) {
	googleValidateErr := s.validateGoogleToken(authBody.ProviderToken)
	if googleValidateErr != nil {
		return nil, nil, googleValidateErr
	}

	user, getUserErr := s.userService.GetUserByEmail(authBody.Email)

	if getUserErr == nil {
		externalLogin := &UpdateExternalLoginBody{
			ClientToken: authBody.ProviderToken,
		}
		_, saveExternalLoginError := s.extLoginService.UpdateUserExternalLoginToken(user.Id, externalLogin)
		if saveExternalLoginError != nil {
			return nil, nil, saveExternalLoginError
		}
		return s.updateTokenPair(user)
	} else {
		userBody := UserBody{
			Email:    authBody.Email,
			UserName: authBody.UserName,
		}
		savedUser, saveUserError := s.userService.SaveUser(userBody)
		if saveUserError != nil {
			return nil, nil, saveUserError
		}

		externalLogin := &models.ExternalLogin{
			ClientId:    authBody.ProviderId,
			ClientToken: authBody.ProviderToken,
			UserRefer:   savedUser.Id,
			Provider:    models.Google,
		}
		_, saveExternalLoginError := s.extLoginService.SaveExternalLogin(externalLogin)
		if saveExternalLoginError != nil {
			// Consider rolling back user creation or logging, for now, return error
			return nil, nil, saveExternalLoginError
		}
		return s.generateAndSaveTokenPair(savedUser)
	}
}

func (s *AuthService) RefreshToken(request RefreshTokenRequest) (*RefreshTokenResponse, error) {
	validationErr := s.AppTokenManager.ValidateToken(request.RefreshToken)
	if validationErr != nil {
		return nil, validationErr
	}

	claims, claimsErr := s.AppTokenManager.GetClaims(request.RefreshToken)
	if claimsErr != nil {
		return nil, claimsErr
	}

	email, ok := claims["email"].(string)
	if !ok {
		return nil, errors.New("email claim is not a string or not found")
	}

	user, getUserErr := s.userService.GetUserByEmail(email)
	if getUserErr != nil {
		return nil, getUserErr
	}

	accessTokenString, accessTokenErr := s.AppTokenManager.GenerateToken(*user, models.Access)
	if accessTokenErr != nil {
		return nil, accessTokenErr
	}

	dbAccessToken, getDbAccessTokenErr := s.tokenService.GetUserTokenByKind(user.Id, models.Access)
	if getDbAccessTokenErr != nil {
		return nil, getDbAccessTokenErr
	}

	accessToken := &models.Token{
		Id:         dbAccessToken.Id,
		TokenValue: accessTokenString,
		Kind:       models.Access,
		UserRefer:  user.Id,
	}

	_, saveAccessTokenErr := s.tokenService.UpdateToken(accessToken)
	if saveAccessTokenErr != nil {
		return nil, saveAccessTokenErr
	}

	return &RefreshTokenResponse{Token: accessToken.TokenValue}, nil
}

func (s *AuthService) generateAndSaveTokenPair(user *models.User) (accessToken *models.Token, refreshToken *models.Token, err error) {
	accessTokenString, accessTokenErr := s.AppTokenManager.GenerateToken(*user, models.Access)
	if accessTokenErr != nil {
		return nil, nil, accessTokenErr
	}
	accessToken = &models.Token{
		TokenValue: accessTokenString,
		Kind:       models.Access,
		UserRefer:  user.Id,
	}

	refreshTokenString, refreshTokenErr := s.AppTokenManager.GenerateToken(*user, models.Refresh)
	if refreshTokenErr != nil {
		return nil, nil, refreshTokenErr
	}
	refreshToken = &models.Token{
		TokenValue: refreshTokenString,
		Kind:       models.Refresh,
		UserRefer:  user.Id,
	}

	_, saveAccessTokenErr := s.tokenService.SaveToken(accessToken)
	if saveAccessTokenErr != nil {
		return nil, nil, saveAccessTokenErr
	}

	_, saveRefreshTokenErr := s.tokenService.SaveToken(refreshToken)
	if saveRefreshTokenErr != nil {
		// Consider cleanup for already saved access token
		return nil, nil, saveRefreshTokenErr
	}
	return accessToken, refreshToken, nil
}

func (s *AuthService) updateTokenPair(user *models.User) (accessToken *models.Token, refreshToken *models.Token, err error) {
	tokenPair, getTokenPairErr := s.tokenService.GetUserTokenPair(user.Id)
	if getTokenPairErr != nil {
		return nil, nil, getTokenPairErr
	}

	accessTokenString, accessTokenErr := s.AppTokenManager.GenerateToken(*user, models.Access)
	if accessTokenErr != nil {
		return nil, nil, accessTokenErr
	}

	refreshTokenString, refreshTokenErr := s.AppTokenManager.GenerateToken(*user, models.Refresh)
	if refreshTokenErr != nil {
		return nil, nil, refreshTokenErr
	}

	var updatedAccess, updatedRefresh *models.Token

	for _, token := range tokenPair {
		var currentTokenToUpdate *models.Token
		if token.Kind == models.Access {
			token.TokenValue = accessTokenString
			updatedAccess = token
			currentTokenToUpdate = updatedAccess
		} else if token.Kind == models.Refresh {
			token.TokenValue = refreshTokenString
			updatedRefresh = token
			currentTokenToUpdate = updatedRefresh
		}

		if currentTokenToUpdate != nil {
			_, updateErr := s.tokenService.UpdateToken(currentTokenToUpdate)
			if updateErr != nil {
				return nil, nil, updateErr // return early on first error
			}
		}
	}

	if updatedAccess == nil || updatedRefresh == nil {
		return nil, nil, errors.New("failed to update token pair, one or both tokens not found in existing pair")
	}

	return updatedAccess, updatedRefresh, nil
}

func (s *AuthService) validateGoogleToken(idToken string) error {
	return s.googleValidator.Validate(idToken)
}

// Interfaces and implementations for the AuthService

// GoogleTokenValidator interface
type GoogleTokenValidator interface {
	Validate(idToken string) error
}

// Interface implementation for GoogleTokenValidator
type DefaultGoogleTokenValidator struct {
	Client           *http.Client
	TokenInfoBaseURL string
}

// Constructor for DefaultGoogleTokenValidator
func NewDefaultGoogleTokenValidator() *DefaultGoogleTokenValidator {
	return &DefaultGoogleTokenValidator{
		TokenInfoBaseURL: "https://www.googleapis.com/oauth2/v3/tokeninfo",
	}
}

// Validate the Google token
func (d *DefaultGoogleTokenValidator) Validate(idToken string) error {
	httpClient := d.Client
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	reqURL := fmt.Sprintf("%s?id_token=%s", d.TokenInfoBaseURL, idToken)
	googleAuthRes, googleAuthReqErr := httpClient.Get(reqURL)
	if googleAuthReqErr != nil {
		return googleAuthReqErr
	}
	defer googleAuthRes.Body.Close()

	if googleAuthRes.StatusCode != http.StatusOK {
		log.Printf("Google token validation failed with status: %s", googleAuthRes.Status)
		return errors.New("google token not valid")
	}
	return nil
}
