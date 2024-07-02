# Extremely simple local monitoring

Plug in your Go http server and use a service such as Uptime Robot to monitor your /status endpoint. Returns http 500 on any failure (plus description in JSON)

```go
import "github.com/gwillem/go-simplemon"

http.HandleFunc("/status", simplemon.Handler)
```
