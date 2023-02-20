package gcp

import (
	"context"
	"errors"
	"fmt"
	"time"

	"google.golang.org/api/idtoken"
)

func ValidateToken(ctx context.Context, token, audience, serviceAccount string) error {
	validator, err := idtoken.NewValidator(ctx)
	if err != nil {
		return fmt.Errorf("failed to create idtoken validator: %w", err)
	}

	payload, err := validator.Validate(ctx, token, audience)
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
