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

package testrig

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/superseriousbusiness/gotosocial/internal/config"
	"github.com/superseriousbusiness/gotosocial/internal/router"
)

// CreateTestContextWithTemplatesAndSessions calls gin.CreateTestContext and then configures the sessions and templates similarly to how the router does
func CreateTestContextWithTemplatesAndSessions(request *http.Request, responseWriter http.ResponseWriter) (*gin.Context, *gin.Engine, sessions.Session) {

	ctx, engine := gin.CreateTestContext(responseWriter)
	ctx.Request = request

	router.LoadTemplateFunctions(engine)

	// does not work because CWD is messed up while tests are running
	// // load templates onto the engine
	// if err := router.LoadTemplates(engine); err != nil {
	// 	panic(err)
	// }

	// https://stackoverflow.com/questions/31873396/is-it-possible-to-get-the-current-root-of-package-structure-as-a-string-in-golan
	_, runtimeCallerLocation, _, _ := runtime.Caller(0)
	projectRoot, err := filepath.Abs(filepath.Join(filepath.Dir(runtimeCallerLocation), "../"))
	if err != nil {
		panic(err)
	}

	templateBaseDir := viper.GetString(config.Keys.WebTemplateBaseDir)

	_, err = os.Stat(filepath.Join(projectRoot, templateBaseDir, "index.tmpl"))
	if err != nil {
		panic(fmt.Errorf("%s doesn't seem to contain the templates; index.tmpl is missing: %s", filepath.Join(projectRoot, templateBaseDir), err))
	}

	tmPath := filepath.Join(projectRoot, fmt.Sprintf("%s*", templateBaseDir))
	engine.LoadHTMLGlob(tmPath)

	store := memstore.NewStore(make([]byte, 32), make([]byte, 32))
	store.Options(router.SessionOptions())

	sessionMiddleware := sessions.Sessions("gotosocial-localhost", store)

	sessionMiddleware(ctx)

	return ctx, engine, sessions.Default(ctx)
}
