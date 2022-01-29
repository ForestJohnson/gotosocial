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

package email

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func loadTemplates(templateBaseDir string) (*template.Template, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("error getting current working directory: %s", err)
	}

	// look for all templates that start with 'email_'
	tmPath := filepath.Join(cwd, fmt.Sprintf("%semail_*", templateBaseDir))
	return template.ParseGlob(tmPath)
}

// See https://pkg.go.dev/net/smtp#SendMail
// for a format example.
func assembleMessage(mailSubject string, mailBody string, mailTo string, mailFrom string) []byte {

	//normalize the mail body to use CRLF line endings
	crlf, _ := regexp.Compile("\\r\\n")
	mailBody = crlf.ReplaceAllString(mailBody, "\n")
	mailBody = strings.ReplaceAll(mailBody, "\n", "\r\n")

	msg := []byte(
		"To: " + mailTo + "\r\n" +
			"Subject: " + mailSubject + "\r\n" +
			"\r\n" +
			mailBody + "\r\n",
	)

	return msg
}
