package context

import (
	"context"

	"github.com/samueldaviddelacruz/go-job-board/API/models"
)

const (
	userKey privateKey = "user"
)

type privateKey string

func WithUser(ctx context.Context, user *models.Company) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func User(ctx context.Context) *models.Company {
	if temp := ctx.Value(userKey); temp != nil {
		if user, ok := temp.(*models.Company); ok {
			return user
		}
	}
	return nil
}
