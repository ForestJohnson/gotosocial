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
	"fmt"
	"net"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/superseriousbusiness/gotosocial/internal/api/model"
	"github.com/superseriousbusiness/gotosocial/internal/config"
	"github.com/superseriousbusiness/gotosocial/internal/db"
	"github.com/superseriousbusiness/gotosocial/internal/gtsmodel"
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
		"reasonRequired": viper.GetBool(config.Keys.AccountsReasonRequired),
	})
}

// RegisterPOSTHandler should be served at https://example.org/auth/register.
// After the user fills out the registration form, including thier email address and desired username/password
// the form will post here, and we will handle thier requested account creation per the configured account settings.
// This handler logic is 95% copy and pasted from the POST /api/v1/accounts handler

func (m *Module) RegisterPOSTHandler(c *gin.Context) {
	l := logrus.WithField("func", "registerPOSTHandler")
	l.Trace("entering registration POST handler")

	s := sessions.Default(c)

	// use the client_id on the session to retrieve info about the app associated with the client_id
	clientID, ok := s.Get(sessionClientID).(string)
	if !ok || clientID == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no client_id found in session"})
		return
	}
	app := &gtsmodel.Application{}
	if err := m.db.GetWhere(c.Request.Context(), []db.Where{{Key: sessionClientID, Value: clientID}}, app); err != nil {
		m.clearSession(s)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("no application found for client id %s", clientID)})
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
	if err := validate.AccountCreationForm(form); err != nil {
		l.Debugf("error validating form: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := m.processor.AccountCreate(c.Request.Context(), app.ID, form)
	if err != nil {
		l.Errorf("internal server error while creating new account: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	l.Tracef("new registered user: %+v", user)

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
