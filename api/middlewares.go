package api

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/adfer-dev/analock-api/auth"
	"github.com/adfer-dev/analock-api/constants"
	"github.com/adfer-dev/analock-api/models"
	"github.com/adfer-dev/analock-api/services"
	"github.com/adfer-dev/analock-api/utils"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
)

var tokenService services.TokenService = &services.DefaultTokenService{}
var userService services.UserService = &services.DefaultUserService{}
var tokenManager auth.TokenManager = auth.NewDefaultTokenManager()
var diaryEntryService services.DiaryEntryService = &services.DefaultDiaryEntryService{}

// AuthMiddleware is a middleware to check if each request is correctly authorized.
// Returs the next http handler to be processed.
func AuthMiddleware(next http.Handler) http.Handler {
	authEndpoints := regexp.MustCompile(constants.ApiV1UrlRoot + `/(auth|swagger)/*`)

	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		//If the endpoint is not allowed, check its auth token.
		if authEndpoints.MatchString(req.URL.Path) {
			next.ServeHTTP(res, req)
		} else {
			authErr := checkAuth(req)

			//If the token is valid, execute the next function. Otherwise, respond with an error.
			if authErr == nil {
				next.ServeHTTP(res, req)
			} else if authErr.Error() != "method not allowed" {
				utils.WriteJSON(res, 401,
					models.HttpError{Status: 401, Description: authErr.Error()})
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

// checkUserOwnershipMiddleware checks if the user has ownership on the resource it is trying to edit or delete.
// Returs the next http handler to be processed.
func checkUserOwnershipMiddleware(req *http.Request) error {
	endpointsToCheck := regexp.MustCompile(
		constants.ApiV1UrlRoot +
			`(` + constants.ApiUrlDiaryEntries +
			`|` + constants.ApiUrlBookRegistrations +
			`|` + constants.ApiUrlGameRegistrations +
			`)/*`)

	if endpointsToCheck.MatchString(req.URL.Path) {
		itemId, _ := strconv.Atoi(mux.Vars(req)["id"])
		tokenValue := req.Header.Get("Authorization")[7:]
		tokenClaims, claimsErr := tokenManager.GetClaims(tokenValue)

		if claimsErr != nil {
			return claimsErr
		}
		tokenEmail := tokenClaims["email"].(string)
		if req.Method == http.MethodGet {
			var ownershipErr error

			if strings.Contains(req.URL.Path, "user") {
				ownershipErr = checkUserEmailOwnership(uint(itemId), tokenEmail)
			} else {
				ownershipErr = checkUserOwnershipFromDiaryEntryId(uint(itemId), tokenEmail)
			}

			if ownershipErr != nil {
				return ownershipErr
			}

		} else if req.Method == http.MethodPut {
			if strings.Contains(req.URL.Path, constants.ApiUrlDiaryEntries) {
				ownershipErr := checkUserOwnershipFromDiaryEntryId(uint(itemId), tokenEmail)

				if ownershipErr != nil {
					return ownershipErr
				}
			}
		}
	}

	return nil
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
	if err := tokenManager.ValidateToken(tokenString); err != nil {
		validationErr, ok := err.(*jwt.ValidationError)
		if ok && validationErr.Errors == jwt.ValidationErrorExpired {
			return errors.New("token expired. Please, get a new one at /auth/refresh-token")
		} else {
			return errors.New("token not valid")
		}
	}

	//Then check if token is in the database
	if _, tokenNotFoundErr := tokenService.GetTokenByValue(tokenString); tokenNotFoundErr != nil {
		return errors.New("token revoked")
	}

	claims, claimsErr := tokenManager.GetClaims(tokenString)

	if claimsErr != nil {
		return claimsErr
	}

	user, _ := userService.GetUserByEmail(claims["email"].(string))
	// user-accessible endpoints
	userDiaryEntryEndpoints := regexp.MustCompile(`/api/v1/diaryEntries/*`)
	userActivityRegistrationEndpoints := regexp.MustCompile(`/api/v1/activityRegistrations/*`)
	// for POST, PUT and DELETE methods, check if user is admin and that the endpoint is not user accessible
	if ((req.Method == "POST" || req.Method == "PUT" || req.Method == "DELETE") && (user.Role != models.Admin)) &&
		(!userDiaryEntryEndpoints.MatchString(req.URL.Path) && !userActivityRegistrationEndpoints.MatchString(req.URL.Path)) {
		return errors.New("method not allowed")
	}

	return nil
}

// Check if a user's email ,identified by the id passed as parameter, corresponds to the email contained in token claims.
func checkUserEmailOwnership(userId uint, tokenEmail string) error {
	entryUser, getUserErr := userService.GetUserById(userId)

	if getUserErr != nil {
		return getUserErr
	}

	if tokenEmail != entryUser.Email {
		return errors.New(constants.ErrorUnauthorizedOperation)
	}
	return nil
}

// Checks if user has ownership of a diary entry, knowing the entry id
func checkUserOwnershipFromDiaryEntryId(itemId uint, tokenEmail string) error {
	diaryEntry, getEntryError := diaryEntryService.GetDiaryEntryById(uint(itemId))

	if getEntryError != nil {
		return getEntryError
	}

	return checkUserEmailOwnership(diaryEntry.Registration.UserRefer, tokenEmail)
}
