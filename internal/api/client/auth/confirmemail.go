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
	"net"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/sirupsen/logrus"
	"github.com/superseriousbusiness/gotosocial/internal/api/model"
	"github.com/superseriousbusiness/gotosocial/internal/oauth"
	"github.com/superseriousbusiness/gotosocial/internal/validate"

	"github.com/gin-gonic/gin"
)

// ConfirmEmailGETHandler should be served at https://example.org/confirm_email. ()
// The idea is to present a registration page to the user, where they can enter their email address, username and password.
// The form will then POST to /confirm_email, which will be handled by RegisterPOSTHandler
func (m *Module) ConfirmEmailGETHandler(c *gin.Context) {
	l := logrus.WithField("func", "confirmEmailGETHandler")

	s := sessions.Default(c)

	// The registration process may be taking place inside an embedded webview (tusky android app) or a different web browser
	// So we aren't guaranteed to be able to access that session from here.
	// But, we can display different helpful text depending on whether the session is present or not.
	userID, hasUserId := s.Get(sessionUserID).(string)

	confirmEmailToken := c.Request.URL.Query().Get("token")

	if confirmEmailToken != "" {

		headerText := "Almost There!"
		paragraphText := "One last step before you can us"
		buttonText := ""

		if m.config.AccountsConfig.RequireApproval {

		}
		c.HTML(http.StatusOK, "email_confirm_landing.tmpl", gin.H{
			"headerText": headerText,
			"headerText": headerText,
			"headerText": headerText,
		})
	} else {
		c.HTML(http.StatusOK, "check_your_email.tmpl", gin.H{})
	}

}

// RegisterPOSTHandler should be served at https://example.org/confirm_email.
// After the user fills out the registration form, including thier email address and desired username/password
// the form will post here, and we will handle thier requested account creation per the configured account settings.
// This handler logic is 95% copy and pasted from the POST /api/v1/accounts handler

func (m *Module) ConfirmEmailPOSTHandler(c *gin.Context) {
	l := logrus.WithField("func", "confirmEmailPOSTHandler")
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

	l.Trace("redirecting to the confirm your email page.")
	c.Redirect(http.StatusFound, ConfirmEmailPath)
}
