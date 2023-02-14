package middleware

import (
	"fmt"
	"net/http"

	"github.com/usa4ev/ghostorange/internal/app/auth/session"
)

// AuthMW returns middleware that enriches the request context with UserID
func AuthorisationMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("Authorization")
		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprint(w, "No Authorization cookie set")
				return
			}

			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Authorization failure: %v", err)

			return
		}

		tokenString := c.Value

		userID, err := session.Verify(tokenString)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "Token verification failure: %v", err)
			return
		}

		next.ServeHTTP(w, session.ReqWithSession(r, userID))
	})
}
