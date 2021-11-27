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
	"net"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/sirupsen/logrus"
	"github.com/superseriousbusiness/gotosocial/internal/api/model"
	"github.com/superseriousbusiness/gotosocial/internal/oauth"
	"github.com/superseriousbusiness/gotosocial/internal/validate"

	"github.com/gin-gonic/gin"
)

// RegisterGETHandler should be served at https://example.org/auth/register.
// The idea is to present a registration page to the user, where they can enter their email address, username and password.
// The form will then POST to /auth/register, which will be handled by RegisterPOSTHandler
func (m *Module) RegisterGETHandler(c *gin.Context) {
	//l := logrus.WithField("func", "RegisterGETHandler")
	if m.idp != nil {
		// TODO
		c.AbortWithError(500, errors.New("TODO what should we do for registration when an Identity Provider (idp) is in use"))
	}
	c.HTML(http.StatusOK, "register.tmpl", gin.H{
		"reasonRequired": m.config.AccountsConfig.ReasonRequired,
	})
}

// RegisterPOSTHandler should be served at https://example.org/auth/register.
// After the user fills out the registration form, including thier email address and desired username/password
// the form will post here, and we will handle thier requested account creation per the configured account settings.
// This handler logic is 95% copy and pasted from the POST /api/v1/accounts handler

func (m *Module) RegisterPOSTHandler(c *gin.Context) {
	l := logrus.WithField("func", "registerPOSTHandler")
	l.Trace("entering registration POST handler")
	authed, err := oauth.Authed(c, true, true, false, false)
	if err != nil {
		l.Debugf("couldn't auth: %s", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	l.Trace("parsing request form")
	form := &model.AccountCreateRequest{}
	if err := c.ShouldBind(form); err != nil || form == nil {
		l.Debugf("could not parse form from request: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing one or more required form values"})
		return
	}

	clientIP := c.ClientIP()
	l.Tracef("attempting to parse client ip address %s", clientIP)
	signUpIP := net.ParseIP(clientIP)
	if signUpIP == nil {
		l.Debugf("error validating sign up ip address %s", clientIP)
		c.JSON(http.StatusBadRequest, gin.H{"error": "ip address could not be parsed from request"})
		return
	}

	form.IP = signUpIP

	l.Tracef("validating form %+v", form)
	if err := validate.AccountCreationForm(form, m.config.AccountsConfig); err != nil {
		l.Debugf("error validating form: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, token, err := m.processor.AccountCreate(c.Request.Context(), authed, form)
	if err != nil {
		l.Errorf("internal server error while creating new account: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	l.Tracef("new registered user: %+v", user)

	// If we are inside an oauth flow right now, it would be nice to finish the flow using this token
	// (but not until after the user has confirmed their email).
	// That way, user can go straight from registering an account to being fully logged in and ready to go.
	// However that would require a lot of extra work and has a lot of edge cases,
	// so for now we will just drop this token on the floor.
	// The user will be instructed to log in again after they confirm thier email.
	l.Tracef("new registered user's token: %+v", token)

	s := sessions.Default(c)
	s.Set(sessionUserID, user.AccountID)
	if err := s.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		m.clearSession(s)
		return
	}

	l.Trace("redirecting to the check your email page.")
	c.Redirect(http.StatusFound, CheckYourEmailPath)
}

// CheckYourEmailGETHandler should be served at https://example.org/check_your_email.
// A page you land on after registering a new account, similar to /setup on mastodon.
func (m *Module) CheckYourEmailGETHandler(c *gin.Context) {
	//l := logrus.WithField("func", "CheckYourEmailGETHandler")
	c.HTML(http.StatusOK, "check_your_email.tmpl", gin.H{"emailFromAddress": m.config.SMTPConfig.From})
}
