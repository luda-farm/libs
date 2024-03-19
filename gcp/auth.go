package gcp

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"google.golang.org/api/idtoken"
)

func ValidateToken(r *http.Request, serviceAccount string) error {
	validator, err := idtoken.NewValidator(r.Context())
	if err != nil {
		return fmt.Errorf("failed to create idtoken validator: %w", err)
	}

	authHeader := r.Header.Get("authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")

	audience := r.URL.String()

	payload, err := validator.Validate(r.Context(), token, audience)
	switch {
	case err != nil:
		return errors.New("invalid token")
	case payload.Expires < time.Now().Unix():
		return errors.New("expired token")
	case payload.Issuer != "https://accounts.google.com":
		return errors.New("invalid token issuer")
	case payload.Claims["email"] != serviceAccount:
		return errors.New("insufficient permissions")
	default:
		return nil
	}
}
