package main

import (
    "fmt"
    "github.com/Danceiny/go.fastjson"
    "github.com/Danceiny/gocelery"
    log "github.com/sirupsen/logrus"
    "os"
    "time"
)

var (
    SITES  map[string]map[string]interface{}
    client *gocelery.CeleryClient
)

func initClient() {
    SITES = make(map[string]map[string]interface{})
    DIANPING := map[string]interface{}{
        "siteName": "dianping",
        "url":      "http://www.dianping.com",
        "subUrls":  []string{"shanghai"},
        // "subUrls2": []string{"ch75", "education",},
        "subUrls2": []string{"ch75/p{uint}"},
    }
    SITES["dianping"] = DIANPING
}

func StartClient() {
    initClient()
    broker := gocelery.NewRedisCeleryBroker(CELERY_BROKER_HOST, CELERY_BROKER_PORT, 0, CELERY_BROKER_PASSWORD)
    backend := gocelery.NewRedisCeleryBackend(CELERY_BACKEND_HOST, CELERY_BACKEND_PORT, 0, CELERY_BACKEND_PASSWORD)
    client, _ = gocelery.NewCeleryClient(broker, backend)
    log.Info("execute crawling task")
    executeTask()
}

func executeTask() {
    var err error
    var task *gocelery.AsyncResult
    task, err = client.DelayKwargs("runTask", SITES["dianping"])
    if err != nil {
        log.Fatal(err)
    }
    result, err := task.Get(24 * time.Hour)
    if err != nil {
        log.Fatal(err)
    }
    log.Infof("%v", result)
    f, err := os.Create(fmt.Sprintf("./client/%d.log", time.Now().Unix()))
    defer f.Close()
    _, _ = f.Write(fastjson.ToJSON(result))
}
