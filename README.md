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

## Example

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

	logrus.Error("This is a ERR log message")
	logrus.Warning("This is a Warning log message")
	logrus.Info("This is a Info log message")
	logrus.Debug("This is a Debug log message")
}
```

This example will generate the following output:

```
Console: 
ERRO[2025-01-11 23:28:10] This is a ERR log message                    
WARN[2025-01-11 23:28:10] This is a Warning log message                
INFO[2025-01-11 23:28:10] This is a Info log message
DEBUG[2025-01-11 23:28:10] This is a Debug log message

WebSockets body example:
{
    "level": "info",
    "msg": "This is a Info log message",
    "time": "2025-01-11T23:18:17+03:00"
}

logs.json:
{"2025-01-11":[
    {"level":"info","msg":"Starting server on :8080","time":"2025-01-11T23:18:31+03:00"},
    {"level":"error","msg":"This is a ERR log message","time":"2025-01-11T23:18:32+03:00"}
  ]
}
