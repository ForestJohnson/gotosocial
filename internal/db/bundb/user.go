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

package bundb

import (
	"context"
	"sync"

	"github.com/ReneKroon/ttlcache"
	"github.com/superseriousbusiness/gotosocial/internal/db"
	"github.com/superseriousbusiness/gotosocial/internal/gtsmodel"
	"github.com/uptrace/bun"
)

type userDB struct {
	conn        *DBConn
	cache       *ttlcache.Cache
	byAccountId *sync.Map
}

func (a *userDB) newUserQ(user *gtsmodel.User) *bun.SelectQuery {
	return a.conn.
		NewSelect().
		Model(user)
}

func (a *userDB) GetUserByUsername(ctx context.Context, accountId string) (*gtsmodel.User, db.Error) {
	user := &gtsmodel.User{}
	err := a.newUserQ(user).Where("user.account_id = ?", accountId).Scan(ctx)
	if err != nil {
		return nil, a.conn.ProcessError(err)
	}
	return user, nil
}
