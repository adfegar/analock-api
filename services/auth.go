package services

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/adfer-dev/analock-api/auth"
	"github.com/adfer-dev/analock-api/models"
)

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

func AuthenticateUser(authBody UserAuthenticateBody) (*models.Token, *models.Token, error) {
	googleValidateErr := validateGoogleToken(authBody.ProviderToken)

	if googleValidateErr != nil {
		return nil, nil, googleValidateErr
	}

	user, getUserErr := GetUserByEmail(authBody.Email)

	if getUserErr == nil {
		externalLogin := &UpdateExternalLoginBody{
			ClientToken: authBody.ProviderToken,
		}
		_, saveExternalLoginError := UpdateUserExternalLoginToken(user.Id, externalLogin)

		if saveExternalLoginError != nil {
			return nil, nil, saveExternalLoginError
		}

		return updateTokenPair(user)
	} else {
		userBody := UserBody{
			Email:    authBody.Email,
			UserName: authBody.UserName,
		}
		savedUser, saveUserError := SaveUser(userBody)

		if saveUserError != nil {
			return nil, nil, saveUserError
		}
		externalLogin := &models.ExternalLogin{
			ClientId:    authBody.ProviderId,
			ClientToken: authBody.ProviderToken,
			UserRefer:   savedUser.Id,
			Provider:    models.Google,
		}
		_, saveExternalLoginError := SaveExternalLogin(externalLogin)

		if saveExternalLoginError != nil {
			return nil, nil, saveExternalLoginError
		}

		return generateAndSaveTokenPair(savedUser)
	}
}

func RefreshToken(request RefreshTokenRequest) (*RefreshTokenResponse, error) {
	validationErr := auth.ValidateToken(request.RefreshToken)

	if validationErr != nil {
		return nil, validationErr
	}

	claims, claimsErr := auth.GetClaims(request.RefreshToken)

	if claimsErr != nil {
		return nil, claimsErr
	}

	user, getUserErr := GetUserByEmail(claims["email"].(string))

	if getUserErr != nil {
		return nil, getUserErr
	}

	accessTokenString, accessTokenErr := auth.GenerateToken(*user, models.Access)

	if accessTokenErr != nil {
		return nil, accessTokenErr
	}

	dbAccessToken, getDbAccessTokenErr := GetUserTokenByKind(user.Id, models.Access)

	if getDbAccessTokenErr != nil {
		return nil, getDbAccessTokenErr
	}

	accessToken := &models.Token{
		Id:         dbAccessToken.Id,
		TokenValue: accessTokenString,
		Kind:       models.Access,
		UserRefer:  user.Id,
	}

	_, saveAccessTokenErr := UpdateToken(accessToken)

	if saveAccessTokenErr != nil {
		return nil, saveAccessTokenErr
	}

	return &RefreshTokenResponse{Token: accessToken.TokenValue}, nil
}

func generateAndSaveTokenPair(user *models.User) (accessToken *models.Token, refreshToken *models.Token, err error) {

	accessTokenString, accessTokenErr := auth.GenerateToken(*user, models.Access)

	if accessTokenErr != nil {
		err = accessTokenErr
		return
	}

	accessToken = &models.Token{
		TokenValue: accessTokenString,
		Kind:       models.Access,
		UserRefer:  user.Id,
	}

	refreshTokenString, refreshTokenErr := auth.GenerateToken(*user, models.Refresh)

	if refreshTokenErr != nil {
		err = refreshTokenErr
		return
	}

	refreshToken = &models.Token{
		TokenValue: refreshTokenString,
		Kind:       models.Refresh,
		UserRefer:  user.Id,
	}

	_, saveAccessTokenErr := SaveToken(accessToken)

	if saveAccessTokenErr != nil {
		err = saveAccessTokenErr
		return
	}

	_, saveRefreshTokenErr := SaveToken(refreshToken)
	if saveRefreshTokenErr != nil {
		err = saveRefreshTokenErr
		return
	}

	return accessToken, refreshToken, err
}

func updateTokenPair(user *models.User) (accessToken *models.Token, refreshToken *models.Token, err error) {
	tokenPair, getTokenPairErr := GetUserTokenPair(user.Id)

	if getTokenPairErr != nil {
		err = getTokenPairErr
		return
	}

	accessTokenString, accessTokenErr := auth.GenerateToken(*user, models.Access)

	if accessTokenErr != nil {
		err = accessTokenErr
		return
	}

	refreshTokenString, refreshTokenErr := auth.GenerateToken(*user, models.Access)

	if refreshTokenErr != nil {
		err = refreshTokenErr
		return
	}

	for _, token := range tokenPair {
		if token.Kind == models.Access {
			token.TokenValue = accessTokenString
			accessToken = token
		} else {
			token.TokenValue = refreshTokenString
			refreshToken = token
		}

		_, updateErr := UpdateToken(token)

		if updateErr != nil {
			err = updateErr
			return
		}
	}

	return accessToken, refreshToken, nil
}

func validateGoogleToken(idToken string) error {
	googleAuthRes, googleAuthReqErr := http.Get(fmt.Sprintf("https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=%s", idToken))

	if googleAuthReqErr != nil {
		return googleAuthReqErr
	}

	if googleAuthRes.StatusCode != 200 {
		log.Println(googleAuthRes)
		return errors.New("google token not valid")
	}

	return nil
}
