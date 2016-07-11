# go-rlogger

Go client for [rlogd](https://github.com/pandax381/rlogd)'s rloggerd.

# Installation

```
go get github.com/hhatto/go-rlogger
```

# Usage

```go
import (
    "github.com/hhatto/go-rlogger"
)

func main() {
    tag := []byte("this.is.tag")
    socketPath := "/var/run/rlogd/rloggerd.sock"
    r := rlogger.NewRLogger(tag, socketPath)
    defer r.Close()

    msg := []byte("Hello rloggerd")
    r.Write(msg)
}
```

# License

MIT
