package main

import (
    "os"
    "strconv"
)

var (
    CELERY_BROKER_HOST      string
    CELERY_BROKER_PORT      int
    CELERY_BROKER_PASSWORD  string
    CELERY_BACKEND_HOST     string
    CELERY_BACKEND_PORT     int
    CELERY_BACKEND_PASSWORD string
)

func init() {
    _ = os.Setenv("CELERY_BROKER_HOST", "127.0.0.1")
    _ = os.Setenv("CELERY_BROKER_PORT", "6379")
    _ = os.Setenv("CELERY_BROKER_PASSWORD", "")
    _ = os.Setenv("CELERY_BACKEND_HOST", "127.0.0.1")
    _ = os.Setenv("CELERY_BACKEND_PORT", "6379")
    _ = os.Setenv("CELERY_BACKEND_PASSWORD", "")

    var tmp string
    var b bool
    if tmp, b = os.LookupEnv("CELERY_BROKER_HOST"); b {
        CELERY_BROKER_HOST = tmp
    }
    if tmp, b = os.LookupEnv("CELERY_BROKER_PORT"); b {
        CELERY_BROKER_PORT, _ = strconv.Atoi(tmp)
    }
    if tmp, b = os.LookupEnv("CELERY_BROKER_PASSWORD"); b {
        CELERY_BROKER_PASSWORD = tmp
    }
    if tmp, b = os.LookupEnv("CELERY_BACKEND_HOST"); b {
        CELERY_BACKEND_HOST = tmp
    }
    if tmp, b = os.LookupEnv("CELERY_BACKEND_PORT"); b {
        CELERY_BACKEND_PORT, _ = strconv.Atoi(tmp)
    }
    if tmp, b = os.LookupEnv("CELERY_BACKEND_PASSWORD"); b {
        CELERY_BACKEND_PASSWORD = tmp
    }
}
