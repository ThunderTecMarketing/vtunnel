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

package main

import (
	_ "github.com/mholt/caddy"
	_ "github.com/mholt/caddy/caddyhttp"
	"github.com/mholt/caddy/caddy/caddymain"
	"flag"
	"github.com/FTwOoO/vtunnel/tunnel"
	_ "github.com/FTwOoO/vtunnel/client"
	_ "github.com/FTwOoO/vtunnel/server"

)

func main() {

	// set caddy file loader,
	// load caddy file as http tunnel type,
	// start the caddy enging
	flag.Set("type", tunnel.ServerType)
	caddymain.Run()
}