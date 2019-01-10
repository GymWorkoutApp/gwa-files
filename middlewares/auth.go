package middlewares

import (
	"encoding/json"
	"fmt"
	"github.com/GymWorkoutApp/gwap-files/errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"log"
	"net/http"
	"os"
	"strings"
)

// BearerAuth parse bearer token
func BearerAuth(r *http.Request) (accessToken string, ok bool) {
	auth := r.Header.Get("Authorization")
	prefix := "Bearer "

	if auth != "" && strings.HasPrefix(auth, prefix) {
		accessToken = auth[len(prefix):]
	} else {
		accessToken = r.FormValue("access_token")
	}

	if accessToken != "" {
		ok = true
	}

	return
}

func MiddlewareAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		accessToken, ok := BearerAuth(c.Request())
		if !ok {
			return errors.NewResponseByError(errors.ErrYouAreNotAuthenticated)
		}

		url := fmt.Sprintf(os.Getenv("INTROSPECT_URL") + "?token=%s", accessToken)

		// Build the request
		resp, err := http.Get(url)
		if err != nil {
			c.Logger().Error(err)
			panic(err)
		}

		defer resp.Body.Close()

		record := map[string]interface{}{
			"access_token": 	 nil,
			"token_type":   	 nil,
			"expires_in":   	 nil,
			"error"		:		 nil,
			"error_description": nil,
		}

		if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
			log.Println(err)
			return errors.NewResponseByError(errors.ErrJsonMarshal)
		}

		if record["error"] != nil {
			response := errors.NewResponseByError(errors.ErrYouAreNotAuthenticated)
			return c.JSON(response.StatusCode, response)
		}

		token, err := jwt.Parse(record["access_token"].(string), func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return []byte(os.Getenv("TOKEN_KEY")), nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			for key, value := range claims {
				c.Set(key, value)
			}
			return next(c)
		} else {
			panic(err)
		}
	}
}