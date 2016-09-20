/*
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 * Author: FTwOoO <booobooob@gmail.com>
 */

package vpn

import (
	"github.com/mholt/caddy/caddyhttp/httpserver"
	"net/http"
)

type handler struct {
	Paths []string
	Next  httpserver.Handler
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	// if the request path is any of the configured paths
	// write hello
	for _, p := range h.Paths {
		if httpserver.Path(r.URL.Path).Matches(p) {
			w.Write([]byte("Hello, I'm a caddy middleware"))
			return 200, nil
		}
	}
	return h.Next.ServeHTTP(w, r)
}
