package router

import (
    "fmt"
    log "github.com/sirupsen/logrus"
    "github.com/tietang/hystrix-go/hystrix"
    "github.com/tietang/props/kvs"
)

const (
    KEY_HYSTRIX_DEFAULT                  = "default"
    KEY_HYSTRIX_TIMEOUT                  = "hystrix.%s.Timeout"
    KEY_HYSTRIX_MAX_CONCURRENT_REQUESTS  = "hystrix.%s.MaxConcurrentRequests"
    KEY_HYSTRIX_REQUEST_VOLUME_THRESHOLD = "hystrix.%s.RequestVolumeThreshold"
    KEY_HYSTRIX_SLEEP_WINDOW             = "hystrix.%s.SleepWindow"
    KEY_HYSTRIX_ERROR_PERCENT_THRESHOLD  = "hystrix.%s.ErrorPercentThreshold"
)
const (
    KEY_CIRCUIT_ENABLED = "hystrix.circuit.enabled"
)

var (
    globalHystrixCommandConfig *hystrix.CommandConfig
)

func getGlobalHystrixCommandConfig(source kvs.ConfigSource) *hystrix.CommandConfig {
    if globalHystrixCommandConfig != nil {
        return globalHystrixCommandConfig
    }

    timeout, err := source.GetInt(fmt.Sprintf(KEY_HYSTRIX_TIMEOUT, KEY_HYSTRIX_DEFAULT))
    if err != nil || timeout <= 0 {
        timeout = hystrix.DefaultTimeout
    }
    maxConcurrentRequests, err := source.GetInt(fmt.Sprintf(KEY_HYSTRIX_MAX_CONCURRENT_REQUESTS, KEY_HYSTRIX_DEFAULT))
    if err != nil || maxConcurrentRequests <= 0 {
        maxConcurrentRequests = hystrix.DefaultMaxConcurrent
    }
    requestVolumeThreshold, err := source.GetInt(fmt.Sprintf(KEY_HYSTRIX_REQUEST_VOLUME_THRESHOLD, KEY_HYSTRIX_DEFAULT))
    if err != nil || requestVolumeThreshold <= 0 {
        requestVolumeThreshold = hystrix.DefaultVolumeThreshold
    }
    sleepWindow, err := source.GetInt(fmt.Sprintf(KEY_HYSTRIX_SLEEP_WINDOW, KEY_HYSTRIX_DEFAULT))
    if err != nil || sleepWindow <= 0 {
        sleepWindow = hystrix.DefaultSleepWindow
    }
    errorPercentThreshold, err := source.GetInt(fmt.Sprintf(KEY_HYSTRIX_ERROR_PERCENT_THRESHOLD, KEY_HYSTRIX_DEFAULT))
    if err != nil || errorPercentThreshold <= 0 {
        errorPercentThreshold = hystrix.DefaultErrorPercentThreshold
    }
    cc := hystrix.CommandConfig{
        Timeout:                timeout,
        MaxConcurrentRequests:  maxConcurrentRequests,
        RequestVolumeThreshold: requestVolumeThreshold,
        SleepWindow:            sleepWindow,
        ErrorPercentThreshold:  errorPercentThreshold,
    }
    log.WithField("GlobalHystrixCommandConfig", cc).Debug()
    return &cc
}

func getHystrixCommandConfig(source kvs.ConfigSource, name string, global *hystrix.CommandConfig) hystrix.CommandConfig {
    timeout, err := source.GetInt(fmt.Sprintf(KEY_HYSTRIX_TIMEOUT, name))
    if err != nil || timeout <= 0 {
        timeout = global.Timeout
    }
    maxConcurrentRequests, err := source.GetInt(fmt.Sprintf(KEY_HYSTRIX_MAX_CONCURRENT_REQUESTS, name))
    if err != nil || maxConcurrentRequests <= 0 {
        maxConcurrentRequests = global.MaxConcurrentRequests
    }
    requestVolumeThreshold, err := source.GetInt(fmt.Sprintf(KEY_HYSTRIX_REQUEST_VOLUME_THRESHOLD, name))
    if err != nil || requestVolumeThreshold <= 0 {
        requestVolumeThreshold = global.RequestVolumeThreshold
    }
    sleepWindow, err := source.GetInt(fmt.Sprintf(KEY_HYSTRIX_SLEEP_WINDOW, name))
    if err != nil || sleepWindow <= 0 {
        sleepWindow = global.SleepWindow
    }
    errorPercentThreshold, err := source.GetInt(fmt.Sprintf(KEY_HYSTRIX_ERROR_PERCENT_THRESHOLD, name))
    if err != nil || errorPercentThreshold <= 0 {
        errorPercentThreshold = global.ErrorPercentThreshold
    }
    cc := hystrix.CommandConfig{
        Timeout:                timeout,
        MaxConcurrentRequests:  maxConcurrentRequests,
        RequestVolumeThreshold: requestVolumeThreshold,
        SleepWindow:            sleepWindow,
        ErrorPercentThreshold:  errorPercentThreshold,
    }
    log.WithField(name+".HystrixCommandConfig", cc).Debug()
    return cc

}

func configHystrix(name string, configSource kvs.ConfigSource) {

    globalHcc := getGlobalHystrixCommandConfig(configSource)
    //配置hystrix
    hcc := getHystrixCommandConfig(configSource, name, globalHcc)
    hystrix.ConfigureCommand(name, hcc)

}
