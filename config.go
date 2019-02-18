package main

import (
    . "github.com/Danceiny/go.utils"
    "os"
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

    CELERY_BROKER_HOST = GetEnvOrDefault("CELERY_BROKER_HOST", "127.0.0.1").(string)
    CELERY_BROKER_PORT = GetEnvOrDefault("CELERY_BROKER_PORT", 6379).(int)
    CELERY_BROKER_PASSWORD = GetEnvOrDefault("CELERY_BROKER_PASSWORD", "").(string)
    CELERY_BACKEND_HOST = GetEnvOrDefault("CELERY_BACKEND_HOST", "127.0.0.1").(string)
    CELERY_BACKEND_PORT = GetEnvOrDefault("CELERY_BACKEND_PORT", 6379).(int)
    CELERY_BACKEND_PASSWORD = GetEnvOrDefault("CELERY_BACKEND_PASSWORD", "").(string)
}
