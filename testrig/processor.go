// GoToSocial
// Copyright (C) GoToSocial Authors admin@gotosocial.org
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package testrig

import (
	"github.com/superseriousbusiness/gotosocial/internal/email"
	"github.com/superseriousbusiness/gotosocial/internal/federation"
	"github.com/superseriousbusiness/gotosocial/internal/media"
	"github.com/superseriousbusiness/gotosocial/internal/processing"
	"github.com/superseriousbusiness/gotosocial/internal/state"
)

// NewTestProcessor returns a Processor suitable for testing purposes
func NewTestProcessor(state *state.State, federator federation.Federator, emailSender email.Sender, mediaManager *media.Manager) *processing.Processor {
	p := processing.NewProcessor(NewTestTypeConverter(state.DB), federator, NewTestOauthServer(state.DB), mediaManager, state, emailSender)
	state.Workers.EnqueueClientAPI = p.EnqueueClientAPI
	state.Workers.EnqueueFederator = p.EnqueueFederator
	return p
}
