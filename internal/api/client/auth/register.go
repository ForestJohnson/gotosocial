/*
   GoToSocial
   Copyright (C) 2021 GoToSocial Authors admin@gotosocial.org

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

package auth

import (
	"errors"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/superseriousbusiness/gotosocial/internal/api/model"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// registration represents a form-submitted new account registration request
type registration struct {
	Email    string `form:"email"`
	Username string `form:"username"`
	Password string `form:"password"`
}

// RegisterGETHandler should be served at https://example.org/auth/register.
// The idea is to present a registration page to the user, where they can enter their email address, username and password.
// The form will then POST to /auth/register, which will be handled by RegisterPOSTHandler
func (m *Module) RegisterGETHandler(c *gin.Context) {
	l := logrus.WithField("func", "RegisterGETHandler")
	l.Trace("entering registration handler")
	if m.idp != nil {
		// TODO
		c.AbortWithError(500, errors.New("TODO what should we do for registration when an Identity Provider (idp) is in use"))
	}
	c.HTML(http.StatusOK, "register.tmpl", gin.H{})
}

// RegisterPOSTHandler should be served at https://example.org/auth/register.
// After the user fills out the registration form, including thier email address and desired username/password
// the form will post here, and we will handle thier requested account creation per the configured account settings
func (m *Module) RegisterPOSTHandler(c *gin.Context) {
	l := logrus.WithField("func", "SignInPOSTHandler")
	if m.idp != nil {
		// TODO
		c.AbortWithError(500, errors.New("TODO what should we do for registration when an Identity Provider (idp) is in use"))
	}

	s := sessions.Default(c)
	form := &model.AccountCreateRequest{}
	if err := c.ShouldBind(form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		m.clearSession(s)
		return
	}
	l.Tracef("parsed registration form: %+v", form)

	s.Set(sessionUserID, userid)
	if err := s.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		m.clearSession(s)
		return
	}

	l.Trace("redirecting to auth page")
	c.Redirect(http.StatusFound, OauthAuthorizePath)
}
