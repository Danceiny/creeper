package main

import (
    "fmt"
    . "github.com/Danceiny/go.utils"
)

// 环境变量
var (
    CELERY_BROKER_HOST      string
    CELERY_BROKER_PORT      int
    CELERY_BROKER_PASSWORD  string
    CELERY_BACKEND_HOST     string
    CELERY_BACKEND_PORT     int
    CELERY_BACKEND_PASSWORD string
    PROXY_SERVER_ADDR       string
    PROXY_SERVER_API        string
    // worker是否进行异步请求
    OPEN_ASYNC_MODE bool
)
// 全局可见变量
var (
    PROXYS []string
)

func init() {
    CELERY_BROKER_HOST = GetEnvOrDefault("CELERY_BROKER_HOST", "127.0.0.1").(string)
    CELERY_BROKER_PORT = GetEnvOrDefault("CELERY_BROKER_PORT", 6379).(int)
    CELERY_BROKER_PASSWORD = GetEnvOrDefault("CELERY_BROKER_PASSWORD", "").(string)
    CELERY_BACKEND_HOST = GetEnvOrDefault("CELERY_BACKEND_HOST", "127.0.0.1").(string)
    CELERY_BACKEND_PORT = GetEnvOrDefault("CELERY_BACKEND_PORT", 6379).(int)
    CELERY_BACKEND_PASSWORD = GetEnvOrDefault("CELERY_BACKEND_PASSWORD", "").(string)
    PROXY_SERVER_ADDR = GetEnvOrDefault("PROXY_SERVER_ADDR", "127.0.0.1:5010").(string)
    PROXY_SERVER_API = fmt.Sprintf("http://%s/get", PROXY_SERVER_ADDR)
    OPEN_ASYNC_MODE = GetEnvOrDefault("ASYNC_MODE", true).(bool)
}
