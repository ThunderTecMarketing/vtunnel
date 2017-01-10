    
A plugin for caddy that implements VPN like wireguard, but use HTTP2 for connection.

## Config file

```
vpn {
    publickey serverpublickey
    privatekey serverprivatekey
    clients {
        publickey client_publickey1
        publickey client_publickey2
        publickey client_publickey3
    }
    mtu 1400
    subnet 192.168.4.1/24
    dnsport 53
    auth /auth
    packet /packet
}
```

