/*
	GoToSocial
	Copyright (C) GoToSocial Authors admin@gotosocial.org
	SPDX-License-Identifier: AGPL-3.0-or-later

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

const React = require("react");

const query = require("../lib/query");

const Loading = require("./loading");

module.exports = function UserLogoutCard() {
	const { data: profile, isLoading } = query.useVerifyCredentialsQuery();
	const { data: instance } = query.useInstanceQuery();
	const [logoutQuery] = query.useLogoutMutation();

	if (isLoading) {
		return <Loading />;
	} else {
		return (
			<div className="account-card">
				<img className="avatar" src={profile.avatar} alt="" />
				<h3 className="text-cutoff">{profile.display_name?.length > 0 ? profile.display_name : profile.acct}</h3>
				<span className="text-cutoff">@{profile.username}@{instance?.account_domain}</span>
				<a onClick={logoutQuery} href="#" aria-label="Log out" title="Log out" className="logout">
					<i className="fa fa-fw fa-sign-out" aria-hidden="true" />
				</a>
			</div>
		);
	}
};