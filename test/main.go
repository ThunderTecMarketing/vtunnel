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
 * key generate with command:
 *  openssl req -new -x509 -days 3650 -newkey rsa:2048 -nodes -keyout 127.0.0.1.key -subj "/C=US/ST=Oregon/L=Portland/CN=127.0.0.1" -out  127.0.0.1.pem
 */

package main

import (
	"flag"
	"fmt"
	"time"
	"net/http"
)

func main() {
	cert := flag.String("cert", "/Users/ganxiangle/Desktop/127.0.0.1.pem", "certificate")
	key := flag.String("key", "/Users/ganxiangle/Desktop/127.0.0.1.key", "private key")
	port := flag.Int("port", 8001, "port")
	flag.Parse()

	err := http.ListenAndServeTLS(fmt.Sprintf(":%d", *port), *cert, *key, &Handler{})
	if err != nil {
		fmt.Println(err)
	}
}

type Handler struct {
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	PrintRequest(r)

	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "%v\n", time.Now())
		w.(http.Flusher).Flush()

	} else if r.Method == http.MethodPost {
		t := r.Header["Content-Type"]
		if t != nil {
			w.Header().Set("Content-Type", t[0])
		}

		buf := make([]byte, 1024)
		n, _ := r.Body.Read(buf)
		w.Header().Set("Content-Length", fmt.Sprintf("%d", n))

		w.Write(buf[:n])
		w.(http.Flusher).Flush()
	}
}

func PrintRequest(r *http.Request) {
	fmt.Printf("%v %s %s\n", r.Proto, r.Method, r.RequestURI)
	fmt.Printf("Body length:%v\n", r.ContentLength)

	for i, x := range r.Header {
		fmt.Printf("%s:%v\n", i, x)
	}
}