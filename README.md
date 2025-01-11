# logstream-go
Duplicating logs via websockets

Introduction
------------

The logstream-go package allows Go programs to conveniently log the program and broadcast logs using WebSockets.
All logs that will be written through the logrus package will be automatically sent to websockets and the log file (if it exists).


## Installation

```bash
$ go get github.com/romanskijdev/logstream-go@latest
```

## Quick Start

Add this import line to the file you're working in:

```Go
import "github.com/romanskijdev/logstream-go"
```

## Init logstream client

```Go
package main

import (
    "fmt"
    "html"
    "net/http"

    "github.com/romanskijdev/logstream-go"
)

func main() {
	path := "logs.json"
	client := logstream.InitLoggerClient(nil, &path)

	client.InitLogger()
	
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", client.HandleConnections)
}
```