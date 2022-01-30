package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/superseriousbusiness/gotosocial/internal/db"
	"github.com/superseriousbusiness/gotosocial/internal/gtserror"
	"github.com/superseriousbusiness/gotosocial/internal/gtsmodel"
	"github.com/superseriousbusiness/oauth2/v4"
)

func (p *processor) OAuthTokenCreate(ctx context.Context, user *gtsmodel.User, oauth2TokenInfo *oauth2.TokenInfo, clientSecret, userid string) (*oauth2.TokenInfo, gtserror.WithCode) {
	if token == "" {
		return nil, gtserror.NewErrorNotFound(errors.New("no token provided"))
	}

	user := &gtsmodel.User{}
	if err := p.db.GetWhere(ctx, []db.Where{{Key: "confirmation_token", Value: token}}, user); err != nil {
		if err == db.ErrNoEntries {
			return nil, gtserror.NewErrorNotFound(err)
		}
		return nil, gtserror.NewErrorInternalError(err)
	}

	if user.Account == nil {
		a, err := p.db.GetAccountByID(ctx, user.AccountID)
		if err != nil {
			return nil, gtserror.NewErrorNotFound(err)
		}
		user.Account = a
	}

	if !user.Account.SuspendedAt.IsZero() {
		return nil, gtserror.NewErrorForbidden(fmt.Errorf("ConfirmEmail: account %s is suspended", user.AccountID))
	}

	if user.UnconfirmedEmail == "" || user.UnconfirmedEmail == user.Email {
		// no pending email confirmations so just return OK
		return user, nil
	}

	if user.ConfirmationSentAt.Before(time.Now().Add(-oneWeek)) {
		return nil, gtserror.NewErrorForbidden(errors.New("ConfirmEmail: confirmation token expired"))
	}

	// mark the user's email address as confirmed + remove the unconfirmed address and the token
	user.Email = user.UnconfirmedEmail
	user.UnconfirmedEmail = ""
	user.ConfirmedAt = time.Now()
	user.ConfirmationToken = ""
	user.UpdatedAt = time.Now()

	if err := p.db.UpdateByPrimaryKey(ctx, user); err != nil {
		return nil, gtserror.NewErrorInternalError(err)
	}

	return m.oauthServer.GenerateUserAccessToken(c.Request.Context(), authed.Token, authed.Application.ClientSecret, user.ID)
}
