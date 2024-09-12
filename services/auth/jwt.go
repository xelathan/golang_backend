package auth

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/xelathan/golang_backend/config"
	"github.com/xelathan/golang_backend/types"
	"github.com/xelathan/golang_backend/utils"
)

type contextKey string

const UserKey contextKey = "userId"

func CreateJWT(secret []byte, userId int) (string, error) {
	expiration := time.Now().Add(time.Second * time.Duration(config.Envs.JWTExpirationInSeconds)).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":    strconv.Itoa(userId),
		"expiredAt": expiration,
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func WithJWTAuth(funcToInvoke http.HandlerFunc, store types.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := getTokenFromRequest(r)
		token, err := validateJWTToken(tokenString)

		if err != nil {
			permissionDenied(w)
			fmt.Printf("%v", err)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			permissionDenied(w)
			return
		}

		str := claims["userId"].(string)

		userId, _ := strconv.Atoi(str)

		user, err := store.GetUserById(userId)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, user.ID)
		r = r.WithContext(ctx)
		funcToInvoke(w, r)
	}
}

func getTokenFromRequest(r *http.Request) string {
	token := r.Header.Get("Authorization")

	if token == "" {
		return ""
	}

	return token
}

func validateJWTToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.Envs.JWTSecret), nil
	})
}

func permissionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
}

func GetUserIdFromContext(ctx context.Context) int {
	userId, ok := ctx.Value(UserKey).(int)

	if !ok {
		return -1
	}

	return userId
}
