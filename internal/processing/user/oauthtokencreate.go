/*
   GoToSocial
   Copyright (C) 2021-2022 GoToSocial Authors admin@gotosocial.org

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package user

import (
	"context"
	"fmt"

	"github.com/superseriousbusiness/gotosocial/internal/gtserror"
	"github.com/superseriousbusiness/gotosocial/internal/gtsmodel"
	"github.com/superseriousbusiness/gotosocial/internal/oauth"
	"github.com/superseriousbusiness/oauth2/v4"
)

func (p *processor) OAuthTokenCreate(ctx context.Context, authed *oauth.Auth, user *gtsmodel.User) (oauth2.TokenInfo, gtserror.WithCode) {
	accessToken, err := p.oauthServer.GenerateUserAccessToken(ctx, authed.Token, authed.Application.ClientSecret, user.ID)
	if err != nil {
		return nil, gtserror.NewErrorInternalError(err, fmt.Sprintf("error generating user access token for user %s", user.ID))
	}
	return accessToken, nil
}
