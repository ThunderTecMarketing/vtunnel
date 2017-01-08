

## Config file

```
realip {
    publickey name
    privatekey   cidr
    clients {
        client_publickey1
        client_publickey2
        client_publickey3
        ...
    }
    
    subnet 192.168.4.1/24
    mtu 1400
    dnsport 53
    authapi /auth
    packetapi /packet
}
```

