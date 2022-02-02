package auth_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"codeberg.org/gruf/go-errors"
	"github.com/stretchr/testify/suite"
	"github.com/superseriousbusiness/gotosocial/internal/api/client/auth"
	"github.com/superseriousbusiness/gotosocial/internal/gtsmodel"
)

type AuthAuthorizeTestSuite struct {
	AuthStandardTestSuite
}

type authorizeHandlerTestCase struct {
	description            string
	mutateUserAccount      func(*gtsmodel.User, *gtsmodel.Account)
	expectedStatusCode     int
	expectedLocationHeader string
}

func (suite *AuthAuthorizeTestSuite) TestAccountAuthorizeHandler() {

	var tests = []authorizeHandlerTestCase{
		{
			description: "user has their email unconfirmed",
			mutateUserAccount: func(user *gtsmodel.User, account *gtsmodel.Account) {
				// nothing to do, weed_lord420 already has their email unconfirmed
			},
			expectedStatusCode:     http.StatusSeeOther,
			expectedLocationHeader: auth.CheckYourEmailPath,
		},
		{
			description: "user has their email confirmed but is not approved",
			mutateUserAccount: func(user *gtsmodel.User, account *gtsmodel.Account) {
				user.ConfirmedAt = time.Now()
				user.Email = user.UnconfirmedEmail
			},
			expectedStatusCode:     http.StatusSeeOther,
			expectedLocationHeader: auth.WaitForApprovalPath,
		},
		{
			description: "user has their email confirmed and is approved, but User entity has been disabled",
			mutateUserAccount: func(user *gtsmodel.User, account *gtsmodel.Account) {
				user.ConfirmedAt = time.Now()
				user.Email = user.UnconfirmedEmail
				user.Approved = true
				user.Disabled = true
			},
			expectedStatusCode:     http.StatusSeeOther,
			expectedLocationHeader: auth.AccountDisabledPath,
		},
		{
			description: "user has their email confirmed and is approved, but Account entity has been suspended",
			mutateUserAccount: func(user *gtsmodel.User, account *gtsmodel.Account) {
				user.ConfirmedAt = time.Now()
				user.Email = user.UnconfirmedEmail
				user.Approved = true
				user.Disabled = false
				account.SuspendedAt = time.Now()
			},
			expectedStatusCode:     http.StatusSeeOther,
			expectedLocationHeader: auth.AccountDisabledPath,
		},
	}

	doTest := func(testCase authorizeHandlerTestCase) {
		ctx, _, testSession := suite.newContext(http.MethodGet, auth.OauthAuthorizePath)

		user := suite.testUsers["unconfirmed_account"]
		account := suite.testAccounts["unconfirmed_account"]

		testSession.Set(sessionUserID, user.ID)
		testSession.Set(sessionClientID, suite.testApplications["application_1"].ClientID)
		if err := testSession.Save(); err != nil {
			panic(errors.WrapMsgf(err, "failed on case: %s", testCase.description))
		}

		testCase.mutateUserAccount(user, account)

		testCase.description = fmt.Sprintf("%s, %t, %s", user.Email, user.Disabled, account.SuspendedAt)

		user.UpdatedAt = time.Now()
		err := suite.db.UpdateByPrimaryKey(context.Background(), user)
		suite.NoError(err)
		_, err = suite.db.UpdateAccount(context.Background(), account)
		suite.NoError(err)

		// call the handler
		suite.authModule.AuthorizeGETHandler(ctx)

		// 1. we should have a redirect
		// suite.Equal(http.StatusSeeOther, recorder.Code)
		suite.Equal(testCase.expectedStatusCode, ctx.Writer.Status(), fmt.Sprintf("failed on case: %s", testCase.description))

		// 2. we should have a redirect to the check your email path, as this user has not confirmed their email yet.
		suite.Equal(testCase.expectedLocationHeader, ctx.Writer.Header().Get("Location"), fmt.Sprintf("failed on case: %s", testCase.description))
	}

	for _, testCase := range tests {
		doTest(testCase)
	}
}

func TestAccountUpdateTestSuite(t *testing.T) {
	suite.Run(t, new(AuthAuthorizeTestSuite))
}
