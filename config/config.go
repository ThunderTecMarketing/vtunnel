/*
 * Author: FTwOoO <booobooob@gmail.com>
 * Created: 2017-03
 */

package config

type ProxyServerSetting struct {
	HTTPVer               string `json:"http_version"`
	BasicProxyCredentials string `json:"proxy_credentials"`
	ProxyAddr             string `json:"addr"`
}

type Config struct {
	ListenAddr  string `json:"addr"`
	LogLevel    string `json:"log_level"`
	ProxyServer ProxyServerSetting `json:"proxy_server"`
}
