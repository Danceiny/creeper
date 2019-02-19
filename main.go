package main

import (
    "github.com/Danceiny/go.fastjson"
)

func main() {
    cmd := fastjson.GetEnvOrDefault("CMD", "worker")
    switch cmd {
    case "worker":
        StartWorker()
        break
    case "client":
        StartClient()
        break
    case "exporter":
        Export()
        break
    }
}
