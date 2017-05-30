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
	"flag"
	_ "github.com/FTwOoO/vtunnel/client"
	_ "github.com/FTwOoO/vtunnel/server"
	"github.com/FTwOoO/vtunnel/tunnel"
	"github.com/mholt/caddy"
	"io/ioutil"
	"log"
	"os"
	"path"
	"github.com/FTwOoO/vtunnel/util"
	"strings"
)

const appName = "vtunnel"
const DEFAULT_CONF = "vtunnel.conf"

var serverType string = tunnel.ServerType

var (
	conf string
)

func init() {

	curDir, err := util.GetCurrentExecDir()
	if err != nil {
		panic(err)
	}
	caddy.DefaultConfigFile = path.Join(curDir, DEFAULT_CONF)

	caddy.TrapSignals()
	flag.StringVar(&conf, "conf", "", "config file to load")
	flag.StringVar(&caddy.PidFile, "pidfile", "", "Path to write pid file")
	flag.BoolVar(&caddy.Quiet, "quiet", false, "Quiet mode (no initialization output)")
	caddy.RegisterCaddyfileLoader("vtunnelLoader", caddy.LoaderFunc(confLoader))

}

func main() {
	flag.Parse()

	caddy.AppName = appName

	var err error

	// Executes Startup events
	caddy.EmitEvent(caddy.StartupEvent)

	// Get Caddyfile input
	caddyfileinput, err := caddy.LoadCaddyfile(serverType)
	if err != nil {
		mustLogFatalf("%v", err)
	}
	// Start your engines
	instance, err := caddy.Start(caddyfileinput)
	if err != nil {
		mustLogFatalf("%v", err)
	}

	// Twiddle your thumbs
	instance.Wait()
}

// mustLogFatalf wraps log.Fatalf() in a way that ensures the
// output is always printed to stderr so the user can see it
// if the user is still there, even if the process log was not
// enabled. If this process is an upgrade, however, and the user
// might not be there anymore, this just logs to the process
// log and exits.
func mustLogFatalf(format string, args ...interface{}) {
	if !caddy.IsUpgrade() {
		log.SetOutput(os.Stderr)
	}
	log.Fatalf(format, args...)
}

func confLoader(serverType string) (caddy.Input, error) {
	if conf == "" {
		return nil, nil
	}

	if conf == "stdin" {
		return caddy.CaddyfileFromPipe(os.Stdin, serverType)
	}

	if !strings.HasPrefix(conf, "/") {
		execDir, _ := util.GetCurrentExecDir()
		conf = path.Join(execDir, conf)
	}

	contents, err := ioutil.ReadFile(conf)
	if err != nil {
		return nil, err
	}
	return caddy.CaddyfileInput{
		Contents:       contents,
		Filepath:       conf,
		ServerTypeName: serverType,
	}, nil
}


