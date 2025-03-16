package api

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/adfer-dev/analock-api/auth"
	"github.com/adfer-dev/analock-api/models"
	"github.com/adfer-dev/analock-api/services"
	"github.com/adfer-dev/analock-api/utils"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
)

// AuthMiddleware is a middleware to check if each request is correctly authorized.
// Returs the next http handler to be processed.
func AuthMiddleware(next http.Handler) http.Handler {
	authEndpoints := regexp.MustCompile(`/api/v1/auth/*`)

	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		//If the endpoint is not allowed, check its auth token.
		if authEndpoints.MatchString(req.URL.Path) {
			next.ServeHTTP(res, req)
		} else {
			authErr := checkAuth(req)

			//If the token is valid, execute the next function. Otherwise, respond with an error.
			if authErr == nil {
				next.ServeHTTP(res, req)
			} else {
				utils.WriteJSON(res, 403,
					models.HttpError{Status: 403, Description: authErr.Error()})
			}
		}
	})
}

// ValidatePathParams checks if the id parameter of an endpoint is a valid number.
// Returs the next http handler to be processed.
func ValidatePathParams(next http.Handler) http.Handler {

	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		idParam, idPresent := mux.Vars(req)["id"]

		//If there is not param, just execute the next function
		if !idPresent {
			next.ServeHTTP(res, req)
		} else {
			if idPresent {
				//If there is param check if it's a number.
				if _, err := strconv.Atoi(idParam); err != nil {
					utils.WriteJSON(res, 400,
						models.HttpError{Status: 400, Description: "Id parameter must be a number."})
				} else {
					next.ServeHTTP(res, req)
				}
			}
		}

	})
}

// CheckUserOwnershipMiddleware checks if the user has ownership on the resource it is trying to edit or delete.
// Returs the next http handler to be processed.
func CheckUserOwnershipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		/*
			if req.Method == "PUT" || req.Method == "DELETE" {

				itemId, _ := strconv.Atoi(mux.Vars(req)["id"])
				tokenValue := req.Header.Get("Authorization")[7:]
				token, err := services.GetTokenByValue(tokenValue)

				if err != nil {
					utils.WriteJSON(res, 500,
						models.HttpError{Status: 500, Description: err.Error()})
				} else {
					//TODO: Implement ownership check
				}
			} else {
				next.ServeHTTP(res, req)
			}
		*/
	})
}

// AUX FUNCTIONS

// checkAuth checks if a request is correctly authorized.
// To a request to be correctly authorized it is needed to provide
// an Authorization header with a valid and unexpired access token.
// Returns error if one of the following happens:
//   - The Authorization header is not provided
//   - The token is expired
//   - The token is not a valid JWT
//   - The request method is not authorized
func checkAuth(req *http.Request) error {
	fullToken := req.Header.Get("Authorization")

	if fullToken == "" || !strings.HasPrefix(fullToken, "Bearer") {
		return errors.New("authorization token must be provided, starting with Bearer")
	}

	tokenString := fullToken[7:]

	//Validate token
	if err := auth.ValidateToken(tokenString); err != nil {
		if err.(*jwt.ValidationError).Errors == jwt.ValidationErrorExpired {
			return errors.New("token expired. Please, get a new one at /auth/refresh-token")
		} else {
			return errors.New("token not valid")
		}
	}

	//Then check if token is in the database
	if _, tokenNotFoundErr := services.GetTokenByValue(tokenString); tokenNotFoundErr != nil {
		return errors.New("token revoked")
	}

	claims, claimsErr := auth.GetClaims(tokenString)

	if claimsErr != nil {
		return claimsErr
	}

	user, _ := services.GetUserByUserName(claims["username"].(string))
	// user-accessible endpoints
	userActionEndpoints := regexp.MustCompile(`/api/v1/users/\d/(add-to-room|delete-from-room)`)
	userNoteEndpoints := regexp.MustCompile(`/api/v1/notes/*`)
	userRoomEndpoints := regexp.MustCompile((`/api/v1/rooms/*`))
	// for POST, PUT and DELETE methods, check if user is admin and that the endpoint is not user accessible
	if ((req.Method == "POST" || req.Method == "PUT" || req.Method == "DELETE") && (user.Role != models.Admin)) &&
		(!userActionEndpoints.MatchString(req.URL.Path) && !userNoteEndpoints.MatchString(req.URL.Path) && !userRoomEndpoints.MatchString(req.URL.Path)) {
		return errors.New("method not allowed")
	}

	return nil
}
