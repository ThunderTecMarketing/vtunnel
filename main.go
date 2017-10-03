package main

import (
	"flag"
	"os"
	"encoding/json"
	"github.com/FTwOoO/vtunnel/config"
	"github.com/FTwOoO/vtunnel/tunnel"
)

var (
	APP_NAME       = "vtunnel"
	flagSet        = flag.NewFlagSet(APP_NAME, flag.ExitOnError)
	configFilePath = flagSet.String("conf", "vtunnel.json", "config file physical path")
)

func NewConfiguration(physicalPath string, configPointer interface{}) {
	file, err := os.Open(physicalPath)
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(configPointer)
	if err != nil {
		panic(err)
	}
}

func ParseFlags() *config.Config {
	flagSet.Parse(os.Args[1:])

	configPointer := &config.Config{}
	NewConfiguration(*configFilePath, configPointer)
	return configPointer
}

func main() {

	c := ParseFlags()

	s, err := tunnel.NewServer(c)
	if err != nil {
		panic(err)
	}

	l, err := s.Listen()
	if err != nil {
		panic(err)
	}

	s.Serve(l)
}
