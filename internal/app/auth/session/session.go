package session

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const secret = "KLjkjd34vnsiullJK23490"

// Verify returns userID and nil as an error if passed token is valid
// and error if invalid
func Verify(signedString string) (string, error) {
	token, err := jwt.Parse(signedString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})

	if err != nil {
		return "", fmt.Errorf("token is not valid: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims["userID"].(string), nil
	}

	return "", fmt.Errorf("token does not contain user id")
}

// Open opens new session and returns
// a signed JWT string with expiration date and UserID
func Open(userID string, lifeTime time.Duration) (string, time.Time, error) {
	expiresAt := time.Now().Add(lifeTime)

	claims := jwt.MapClaims{
		"userID": userID,
		"exp":    expiresAt.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedString, err := token.SignedString([]byte(secret))

	return signedString, expiresAt, err
}


const ctxKeyUserID contextKey = 0 // key to a userID context value

var sessionErrNoUserID = errors.New("request ctx does not contain userID key")

type contextKey int

// ReqWithSession adds userID value with CtxKeyUserID key 
// to a given ctx
func ReqWithSession(r *http.Request, usrID string) *http.Request {
	ctx := context.WithValue(r.Context(), ctxKeyUserID, usrID)
	return r.WithContext(ctx)
}

// UserIDFromCtx return userID from reques ctx value.
// Throws err if theres none
func UserIDFromCtx(r *http.Request) (string, error){
	if val := r.Context().Value(ctxKeyUserID); val == nil{
		return "", sessionErrNoUserID
	}else{
		return val.(string), nil
	}
}