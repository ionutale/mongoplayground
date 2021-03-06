// mongoplayground: a sandbox to test and share MongoDB queries
// Copyright (C) 2017 Adrien Petel
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestView(t *testing.T) {

	defer testServer.clearDatabases(t)

	viewTests := []struct {
		name         string
		params       url.Values
		url          string
		responseCode int
	}{
		{
			name:         "template parameters",
			params:       templateParams,
			url:          templateURL,
			responseCode: http.StatusOK,
		},
		{
			name: "new config",
			params: url.Values{
				"mode":   {"bson"},
				"config": {`[{"_id": 1}]`},
				"query":  {templateQuery},
			},
			url:          "p/DEz-pkpheLX",
			responseCode: http.StatusOK,
		},
		{
			name:         "non existing url",
			params:       templateParams,
			url:          "p/unknownURL",
			responseCode: http.StatusNotFound,
		},
		{
			name:         "url with extra param",
			params:       templateParams,
			url:          templateURL + "&usg",
			responseCode: http.StatusOK,
		},
		{
			name:         "url with invalid id length",
			params:       templateParams,
			url:          "p/short",
			responseCode: http.StatusNotFound,
		},
	}

	// start by saving all needed playground
	for _, tt := range viewTests {
		httpBody(t, testServer.saveHandler, http.MethodPost, saveEndpoint, tt.params)
	}

	t.Run("parallel view", func(t *testing.T) {
		for _, tt := range viewTests {

			test := tt // capture range variable
			t.Run(test.name, func(t *testing.T) {

				t.Parallel()

				req, _ := http.NewRequest(http.MethodGet, "/"+test.url, nil)
				resp := httptest.NewRecorder()
				testServer.viewHandler(resp, req)

				if test.responseCode != resp.Code {
					t.Errorf("expected response code %d but got %d", test.responseCode, resp.Code)
				}
			})
		}
	})
}
