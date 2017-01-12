# Centrifugo Benchmark

Benchmark for centrifugo

## Config benchmark

Fix credentials and ws url.
```
secret := "admin"
wsURL := "ws://ws.bagoo.io/connection/websocket"
```

Edit nbClient and nbMessag.
```
const nbClient = 1000
const nbMessage = 100
```

Edit speed.
```
time.Sleep(70 * time.Millisecond)
```
