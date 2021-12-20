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

package nodeinfo

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/superseriousbusiness/gotosocial/internal/api"
)

// NodeInfoGETHandler swagger:operation GET /nodeinfo/2.0 nodeInfoGet
//
// Returns a compliant nodeinfo response to node info queries.
//
// See: https://nodeinfo.diaspora.software/schema.html
//
// ---
// tags:
// - nodeinfo
//
// produces:
// - application/json; profile="http://nodeinfo.diaspora.software/ns/schema/2.0#"
//
// responses:
//   '200':
//     schema:
//       "$ref": "#/definitions/nodeinfo"
func (m *Module) NodeInfoGETHandler(c *gin.Context) {
	l := logrus.WithFields(logrus.Fields{
		"func":       "NodeInfoGETHandler",
		"user-agent": c.Request.UserAgent(),
	})

	if _, err := api.NegotiateAccept(c, api.JSONAcceptHeaders...); err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}

	ni, err := m.processor.GetNodeInfo(c.Request.Context(), c.Request)
	if err != nil {
		l.Debugf("error with get node info request: %s", err)
		c.JSON(err.Code(), err.Safe())
		return
	}

	b, jsonErr := json.Marshal(ni)
	if jsonErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": jsonErr.Error()})
	}

	c.Data(http.StatusOK, `application/json; profile="http://nodeinfo.diaspora.software/ns/schema/2.0#"`, b)
}
