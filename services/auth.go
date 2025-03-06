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
	ProviderId    string `json:"providerId" validate:"required"`
	ProviderToken string `json:"providerToken" validate:"required,jwt"`
}

type UserRegisterBody struct {
	UserName      string `json:"username" validate:"required"`
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

func RegisterUser(userBody UserRegisterBody) (*models.Token, *models.Token, error) {
	user, saveUserErr := SaveUser(UserBody{userBody.UserName})

	if saveUserErr != nil {
		return nil, nil, saveUserErr
	}

	googleTokenErr := validateGoogleToken(userBody.ProviderToken)

	if googleTokenErr != nil {
		return nil, nil, googleTokenErr
	}

	_, saveExternalLoginErr := SaveExternalLogin(&models.ExternalLogin{Provider: models.Google,
		ClientId: userBody.ProviderId, ClientToken: userBody.ProviderToken, UserRefer: user.Id})

	if saveExternalLoginErr != nil {
		return nil, nil, saveExternalLoginErr
	}

	return generateAndSaveTokenPair(user)
}

func AuthenticateUser(authBody UserAuthenticateBody) (*models.Token, *models.Token, error) {
	googleValidateErr := validateGoogleToken(authBody.ProviderToken)

	if googleValidateErr != nil {
		return nil, nil, googleValidateErr
	}

	externalLogin, err := GetExternalLoginByClientId(authBody.ProviderId)

	if err != nil {
		return nil, nil, err
	}

	user, getUserErr := GetUserById(externalLogin.UserRefer)

	if getUserErr != nil {
		return nil, nil, getUserErr
	}

	return updateTokenPair(user)
}

func RefreshToken(request RefreshTokenRequest) (*models.Token, error) {
	validationErr := auth.ValidateToken(request.RefreshToken)

	if validationErr != nil {
		return nil, validationErr
	}

	claims, claimsErr := auth.GetClaims(request.RefreshToken)

	if claimsErr != nil {
		return nil, claimsErr
	}

	user, getUserErr := GetUserByUserName(claims["username"].(string))

	if getUserErr != nil {
		return nil, getUserErr
	}

	accessTokenString, accessTokenErr := auth.GenerateToken(*user, models.Access)

	if accessTokenErr != nil {
		return nil, accessTokenErr
	}

	accessToken := &models.Token{
		TokenValue: accessTokenString,
		Kind:       models.Access,
		UserRefer:  user.Id,
	}

	_, saveAccessTokenErr := SaveToken(accessToken)

	if saveAccessTokenErr != nil {
		return nil, saveAccessTokenErr
	}

	return accessToken, nil
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
