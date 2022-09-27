package httputil

import (
	"net/http"
	"strings"
	"time"

	"github.com/luda-farm/libs/std"
	"google.golang.org/api/idtoken"
)

func ValidateGoogleToken(audience, serviceAccount string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix(r.Header.Get("authorization"), "Bearer ")
		validator := std.Must(idtoken.NewValidator(r.Context()))
		payload, err := validator.Validate(r.Context(), token, audience)
		switch {
		case err != nil:
			http.Error(w, "`invalid token`", http.StatusUnauthorized)
		case payload.Expires < time.Now().Unix():
			http.Error(w, "`expired token`", http.StatusUnauthorized)
		case payload.Issuer != "https://accounts.google.com":
			http.Error(w, "`invalid issuer`", http.StatusUnauthorized)
		case payload.Claims["email"] != serviceAccount:
			http.Error(w, "`insufficient permissions`", http.StatusForbidden)
		default:
			next(w, r)
		}
	}
}
