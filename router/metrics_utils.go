package router

import (
	"fmt"
	"github.com/rcrowley/go-metrics"
	"github.com/tietang/props/kvs"
	"github.com/tietang/zebra/meter"
	"time"
)

const (
	METER_ERROR_REQUEST_PREFIX = "error:"
	METER_OK_REQUEST_PREFIX    = "ok:"
	METER_50X_REQUEST_PREFIX   = "50x:"
	METER_40X_REQUEST_PREFIX   = "40x:"

	KEY_LB_NAME_TEMPLATE          = "%s.balancer.name"
	KEY_FIBONACCI_BASE_TEMPLATE   = "%s.fibonacci.base"
	KEY_MAX_FAILS_TEMPLATE        = "%s.max.fails"
	KEY_FAIL_TIME_WINDOW_TEMPLATE = "%s.fail.time.window"
	KEY_FAIL_SLEEP_MODE_TEMPLATE  = "%s.fail.sleep.mode"
	KEY_FAIL_SLEEP_X_TEMPLATE     = "%s.fail.sleep.x"
	KEY_FAIL_SLEEP_MAX_TEMPLATE   = "%s.fail.sleep.max"
)

var UrlRegistry metrics.Registry = meter.NewRegistry()
var ServiceRegistry metrics.Registry = meter.NewRegistry()
var InstanceRegistry metrics.Registry = meter.NewRegistry()

func metricsErrorKey(ins *HostInstance) string {
	return METER_ERROR_REQUEST_PREFIX + ins.Name
}

func GetOrRegisterErrorMeter(ins *HostInstance, seconds int) meter.MeterX {
	key := metricsErrorKey(ins)
	m := meter.GetOrRegisterMeterX(key, seconds, InstanceRegistry)
	return m
}

func GetOrRegisterErrorMeterSnapshot(ins *HostInstance, seconds int) meter.MeterX {
	key := metricsErrorKey(ins)
	m := meter.GetOrRegisterMeterX(key, seconds, InstanceRegistry)
	return m.Snapshot()
}

func getAppFailSleepSeqX(conf kvs.ConfigSource, name string) []int {
	defaultx := conf.Ints(fmt.Sprintf(KEY_FAIL_SLEEP_X_TEMPLATE, LB_DEFAULT_NAME))
	if len(defaultx) == 0 {
		defaultx = DEFAULT_FAIL_SLEEP_SEQ_X
	}
	x := conf.Ints(fmt.Sprintf(KEY_FAIL_SLEEP_X_TEMPLATE, name))
	if len(x) == 0 {
		x = defaultx
	}
	return x
}

//熔断窗口倍数
func getAppFailSleepX(conf kvs.ConfigSource, name string) int {
	defaultx := conf.GetIntDefault(fmt.Sprintf(KEY_FAIL_SLEEP_X_TEMPLATE, LB_DEFAULT_NAME), DEFAULT_FAIL_SLEEP_X)
	x := conf.GetIntDefault(fmt.Sprintf(KEY_FAIL_SLEEP_X_TEMPLATE, name), defaultx)
	return x
}

//配置的失败时间窗口模式
func getAppFailSleepMode(conf kvs.ConfigSource, name string) string {
	defaultMode := conf.GetDefault(fmt.Sprintf(KEY_FAIL_SLEEP_MODE_TEMPLATE, LB_DEFAULT_NAME), DEFAULT_FAIL_SLEEP_MODE)
	mode := conf.GetDefault(fmt.Sprintf(KEY_FAIL_SLEEP_MODE_TEMPLATE, name), defaultMode)
	return mode
}

//配置的失败时间窗口
func getAppFailTimeWindow(conf kvs.ConfigSource, name string) time.Duration {
	defaultMode := conf.GetDurationDefault(fmt.Sprintf(KEY_FAIL_TIME_WINDOW_TEMPLATE, LB_DEFAULT_NAME), DEFAULT_FAIL_TIME_WINDOW)
	mode := conf.GetDurationDefault(fmt.Sprintf(KEY_FAIL_TIME_WINDOW_TEMPLATE, name), defaultMode)
	return mode
}

//配置的失败时间窗口，秒数
func getAppFailTimeWindowSeconds(conf kvs.ConfigSource, name string) int {
	d := getAppFailTimeWindow(conf, name)
	return int(d.Nanoseconds() / time.Second.Nanoseconds())
}

//配置的时间窗口内最大失败次数
func getAppMaxFails(conf kvs.ConfigSource, name string) int {
	defaultFails := conf.GetIntDefault(fmt.Sprintf(KEY_MAX_FAILS_TEMPLATE, LB_DEFAULT_NAME), DEFAULT_LB_MAX_FAILS)
	fails := conf.GetIntDefault(fmt.Sprintf(KEY_MAX_FAILS_TEMPLATE, name), defaultFails)
	return fails

}

//配置的最大sleep时间
func getAppMaxFailedSleep(conf kvs.ConfigSource, name string) time.Duration {
	defaultTime := conf.GetDurationDefault(fmt.Sprintf(KEY_FAIL_SLEEP_MAX_TEMPLATE, LB_DEFAULT_NAME), DEFAULT_FAIL_SLEEP_MAX)
	fails := conf.GetDurationDefault(fmt.Sprintf(KEY_FAIL_SLEEP_MAX_TEMPLATE, name), defaultTime)
	return fails
}

//配置的最大sleep时间秒数
func getAppMaxFailedSleepSeconds(conf kvs.ConfigSource, name string) int {
	d := getAppMaxFailedSleep(conf, name)
	return int(d.Nanoseconds() / time.Second.Nanoseconds())
}
