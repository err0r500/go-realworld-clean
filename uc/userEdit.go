package uc

import (
	"context"

	"github.com/err0r500/go-realworld-clean/domain"
	"github.com/opentracing/opentracing-go"
)

func (i interactor) UserEdit(ctx context.Context, userName string, fieldsToUpdate map[domain.UserUpdatableProperty]*string) (*domain.User, string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "uc:user_edit")
	defer span.Finish()

	user, ok := i.userRW.GetByName(ctx, userName)
	if !ok {
		return nil, "", ErrTechnical
	}
	if user == nil {
		return nil, "", ErrNotFound
	}
	if user.Name != userName {
		return nil, "", ErrUnauthorized
	}

	domain.UpdateUser(user,
		domain.SetUserName(fieldsToUpdate[domain.UserName]),
		domain.SetUserEmail(fieldsToUpdate[domain.UserEmail]),
		domain.SetUserBio(fieldsToUpdate[domain.UserBio]),
		domain.SetUserImageLink(fieldsToUpdate[domain.UserImageLink]),
		domain.SetUserPassword(fieldsToUpdate[domain.UserPassword]),
	)

	if err := i.userValidator.CheckUser(*user); err != nil {
		return nil, "", err
	}

	if ok := i.userRW.Save(ctx, *user); !ok {
		return nil, "", ErrTechnical
	}

	token, err := i.authHandler.GenUserToken(user.Name)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}
