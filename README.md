    
TCP/UDP tunnel like [GitHub - shell909090/goproxy](https://github.com/shell909090/goproxy). 

## Features
* Base on caddy framework with new server type.
* SOCKS5 for local interface
* Layers of encoders including msgpack. 
* Multiplexing for connections

## Config

### Server
The server listen on `ftwo.me:10809`, using the transport `TCP-Fragment-AheadGCM256-Msgpack` and transport
key `e01ee3207ea15d346c362b7e20cef3a1088ec0a11a1141b3584ed44e2bb69531`:

```
0.0.0.0:10809 {
    server {
       transportType TCP-Fragment-AheadGCM256-Msgpack
       transportKey  e01ee3207ea15d346c362b7e20cef3a1088ec0a11a1141b3584ed44e2bb69531

       logFile "./vtunnel_server.log"
       logLevel DEBUG
    }
}
```

### Client
The client listen on `localhost:1080` as SOCKS5 server, and forward connections to the vtunnel server `ftwo.me:10809`:

```
localhost:1080 {
    client {
       proxyType  socks5
       remoteAddr  ftwo.me:10809

       transportType TCP-Fragment-AheadGCM256-Msgpack
       transportKey  e01ee3207ea15d346c362b7e20cef3a1088ec0a11a1141b3584ed44e2bb69531

       logFile "./vtunnel_client.log"
       logLevel DEBUG
    }
}
```
