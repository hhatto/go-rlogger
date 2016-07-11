# go-rlogger [![](https://travis-ci.org/hhatto/go-rlogger.svg?branch=master)](https://travis-ci.org/hhatto/go-rlogger)

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
    socketPath := "/var/run/rlogd/rloggerd.sock"
    r := rlogger.NewRLogger(socketPath)
    defer r.Close()

    tag := []byte("this.is.tag")
    msg := []byte("Hello rloggerd")
    r.Write(tag, msg)
}
```

# X-rlogger family
* [rlogger-py](https://github.com/KLab/rlogger-py) (Python)
* [php-rlogger](https://github.com/hnw/php-rlogger) (PHP)

# License

MIT
