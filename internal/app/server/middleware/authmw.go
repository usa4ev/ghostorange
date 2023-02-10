package middleware

import (
	"net/http"

	"ghostorange/internal/app/auth/session"
)

// AuthMW returns middleware that enriches the request context with UserID
func AuthorisationMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("Authorization")
		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)

				return
			}

			w.WriteHeader(http.StatusBadRequest)

			return
		}

		tokenString := c.Value

		userID, err := session.Verify(tokenString)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)

			return
		}

		next.ServeHTTP(w, session.ReqWithSession(r, userID))
	})
}



