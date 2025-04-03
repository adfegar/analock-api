package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/adfer-dev/analock-api/auth"
	"github.com/adfer-dev/analock-api/models"
	"github.com/adfer-dev/analock-api/services"
	"github.com/adfer-dev/analock-api/utils"
	"github.com/gorilla/mux"
)

func InitAuthRoutes(router *mux.Router) {
	router.HandleFunc("/api/v1/auth/authenticate", utils.ParseToHandlerFunc(handleAuthenticateUser)).Methods("POST")
	router.HandleFunc("/api/v1/auth/refreshToken", utils.ParseToHandlerFunc(handleRefreshToken)).Methods("POST")
}

func handleAuthenticateUser(res http.ResponseWriter, req *http.Request) error {
	authenticateBody := services.UserAuthenticateBody{}

	validationErrs := utils.HandleValidation(req, &authenticateBody)

	if len(validationErrs) > 0 {
		return utils.WriteJSON(res, 400, validationErrs)
	}

	accessToken, refreshToken, authErr := services.AuthenticateUser(authenticateBody)

	if authErr != nil {
		return utils.WriteJSON(res, 500, models.HttpError{Status: http.StatusInternalServerError, Description: "Error happenned when authenticating user. Please, try again."})
	}

	claims, claimsErr := auth.GetClaims(refreshToken.TokenValue)

	if claimsErr != nil {
		return claimsErr
	}
	res.Header().Add("Set-Cookie", fmt.Sprintf("refreshToken=%s; Expires=%d; HttpOnly", refreshToken.TokenValue, int64(claims["exp"].(float64))))
	return utils.WriteJSON(res, 200,
		services.TokenResponse{AccessToken: accessToken.TokenValue, RefreshToken: refreshToken.TokenValue})
}

func handleRefreshToken(res http.ResponseWriter, req *http.Request) error {
	authenticateBody := services.RefreshTokenRequest{}

	validationErrs := utils.HandleValidation(req, &authenticateBody)

	if len(validationErrs) > 0 {
		return utils.WriteJSON(res, 403, validationErrs)
	}

	newAccessToken, refreshTokenErr := services.RefreshToken(authenticateBody)

	log.Println(refreshTokenErr)

	if refreshTokenErr != nil {
		return utils.WriteJSON(res, 403, refreshTokenErr)
	}

	return utils.WriteJSON(res, 200, newAccessToken)
}
