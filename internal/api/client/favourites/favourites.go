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

package favourites

import (
	"net/http"

	"github.com/superseriousbusiness/gotosocial/internal/api"
	"github.com/superseriousbusiness/gotosocial/internal/processing"
	"github.com/superseriousbusiness/gotosocial/internal/router"
)

const (
	// BasePath is the base URI path for serving favourites
	BasePath = "/api/v1/favourites"

	// MaxIDKey is the url query for setting a max status ID to return
	MaxIDKey = "max_id"
	// SinceIDKey is the url query for returning results newer than the given ID
	SinceIDKey = "since_id"
	// MinIDKey is the url query for returning results immediately newer than the given ID
	MinIDKey = "min_id"
	// LimitKey is for specifying maximum number of results to return.
	LimitKey = "limit"
	// LocalKey is for specifying whether only local statuses should be returned
	LocalKey = "local"
)

// Module implements the ClientAPIModule interface for everything relating to viewing favourites
type Module struct {
	processor processing.Processor
}

// New returns a new favourites module
func New(processor processing.Processor) api.ClientModule {
	return &Module{
		processor: processor,
	}
}

// Route attaches all routes from this module to the given router
func (m *Module) Route(r router.Router) error {
	r.AttachHandler(http.MethodGet, BasePath, m.FavouritesGETHandler)
	return nil
}
