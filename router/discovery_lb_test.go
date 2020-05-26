package router

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tietang/go-eureka-client/eureka"
	"github.com/tietang/props/kvs"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestDiscoveryBalancer_AddHostInstance(t *testing.T) {
	conf := kvs.NewEmptyMapConfigSource("map")
	d := NewDiscoveryBalancer(conf)
	defer d.stop()
	appName := "test1"
	address, port := "192.168.1.2", "8080"
	h := newHostInstance(appName, address, port)
	d.AddHostInstance(appName, h)
	address, port = "192.168.1.3", "8080"
	h = newHostInstance(appName, address, port)
	d.AddHostInstance(appName, h)
	//
	appName = "test2"
	address, port = "192.168.1.2", "8081"
	h = newHostInstance(appName, address, port)
	d.AddHostInstance(appName, h)
	address, port = "192.168.1.3", "8081"
	h = newHostInstance(appName, address, port)
	d.AddHostInstance(appName, h)
	//
	Convey("all", t, func() {
		//So(len(d.Hosts.s), ShouldEqual, 2)
		//So(len(d.UnavailableHosts), ShouldEqual, 0)
		So(d.Balancer, ShouldNotBeNil)
		Convey("add 1 instance", func() {
			ins := d.Next(appName, appName, false)
			So(ins, ShouldNotBeNil)
			ins = d.Next(appName, appName, false)
			So(ins, ShouldNotBeNil)
			ins = d.Next(appName, appName, false)
			So(ins, ShouldNotBeNil)
			ins = d.Next(appName, appName, false)
			So(ins, ShouldNotBeNil)
		})
		Convey("update 1 instance", func() {

			appName = "test2"
			address, port = "192.168.1.2", "8081"
			h = newHostInstance(appName, address, port)
			h.Status = eureka.StatusUp
			d.AddHostInstance(appName, h)

			ins := d.Next(appName, appName, false)
			So(ins, ShouldNotBeNil)
			ins = d.Next(appName, appName, false)
			So(ins, ShouldNotBeNil)
			ins = d.Next(appName, appName, false)
			So(ins, ShouldNotBeNil)
			ins = d.Next(appName, appName, false)
			So(ins, ShouldNotBeNil)
		})

		Convey("add 1 un instance", func() {

			appName = "test2"
			address, port = "192.168.1.2", "8081"
			h = newHostInstance(appName, address, port)
			h.Status = eureka.StatusUp
			d.AddHostInstance(appName, h)
			d.AddUnavailableHostInstance(appName, h)

			ins := d.Next(appName, appName, false)
			So(ins, ShouldNotBeNil)
			ins = d.Next(appName, appName, false)
			So(ins, ShouldNotBeNil)
			ins = d.Next(appName, appName, false)
			So(ins, ShouldNotBeNil)
			ins = d.Next(appName, appName, false)
			So(ins, ShouldNotBeNil)
		})

	})
}

func TestDiscoveryBalancer_NextHostInstance(t *testing.T) {
	conf := kvs.NewEmptyMapConfigSource("map")

	Convey("test nextHostInstance", t, func() {

		Convey(" all un instance", func() {
			d := NewDiscoveryBalancer(conf)
			defer d.stop()
			appName := "test2"
			address, port := "192.168.1.2", "8081"
			h := newHostInstance(appName, address, port)
			h.Status = eureka.StatusUp
			d.AddHostInstance(appName, h)
			d.AddUnavailableHostInstance(appName, h)

			ins := d.Next(appName, appName, false)
			So(ins, ShouldNotBeNil)
			ins = d.Next(appName, appName, false)
			So(ins, ShouldNotBeNil)
			ins = d.Next(appName, appName, false)
			So(ins, ShouldNotBeNil)
			ins = d.Next(appName, appName, false)
			So(ins, ShouldNotBeNil)
		})
		Convey(" all instance is down", func() {
			d := NewDiscoveryBalancer(conf)
			defer d.stop()
			appName := "test2"
			address, port := "192.168.1.2", "8081"
			h := newHostInstance(appName, address, port)
			h.Status = STATUS_DOWN
			d.AddHostInstance(appName, h)
			d.AddUnavailableHostInstance(appName, h)
			//
			address, port = "192.168.1.3", "8081"
			h = newHostInstance(appName, address, port)
			h.Status = STATUS_DOWN
			d.AddHostInstance(appName, h)
			d.AddUnavailableHostInstance(appName, h)

			ins := d.Next(appName, appName, false)
			So(ins, ShouldBeNil)
			ins = d.Next(appName, appName, false)
			So(ins, ShouldBeNil)
			ins = d.Next(appName, appName, false)
			So(ins, ShouldBeNil)
			ins = d.Next(appName, appName, false)
			So(ins, ShouldBeNil)
			ins = d.Next(appName, appName, false)
			So(ins, ShouldBeNil)
			ins = d.Next(appName, appName, false)
			So(ins, ShouldBeNil)
			ins = d.Next(appName, appName, false)
			So(ins, ShouldBeNil)
		})

	})
}

func newHostInstance(appName, address, port string) *HostInstance {
	h := &HostInstance{
		AppGroupName: "DefaultGroup",
		AppName:      appName,
		Address:      address,
		Port:         port,
		InstanceId:   strings.Join([]string{appName, address, port}, ":"),
		Name:         strings.Join([]string{address, port}, ":"),
		Status:       STATUS_UP,
	}
	h.Init()
	return h
}

func TestDiscoveryBalancer_IsOpenFailedSleep_ForFixed(t *testing.T) {
	DiscoveryBalancer_IsOpenFailedSleep_ForFixed_3x(t, 1, 1)
	DiscoveryBalancer_IsOpenFailedSleep_ForFixed_3x(t, 3, 2)
}

func DiscoveryBalancer_IsOpenFailedSleep_ForFixed_3x(t *testing.T, n, sleepx int) {
	conf := kvs.NewEmptyMapConfigSource("map")
	conf.Set("lb.default.balancer.name", "WeightRobinRound")
	conf.Set("lb.default.max.fails", "3")
	conf.Set("lb.default.fail.time.window", strconv.Itoa(n)+"s")
	conf.Set("lb.default.fail.sleep.mode", "fixed")
	conf.Set("lb.default.fail.sleep.x", strconv.Itoa(sleepx))
	conf.Set("lb.default.fail.sleep.max", "6s")
	d := NewDiscoveryBalancer(conf)
	//
	appName := "testForOpenFailedSleep"

	address, port := "192.168.1.2", "8081"
	h1 := newHostInstance(appName, address, port)
	d.AddHostInstance(appName, h1)
	//
	seconds := getAppFailTimeWindowSeconds(conf, h1.AppName)
	m := GetOrRegisterErrorMeter(h1, seconds)
	x := int64(getAppFailSleepX(conf, h1.AppName))
	sleep := getAppMaxFailedSleep(conf, h1.AppName)
	//
	Convey("test IsOpenFailedSleep "+strconv.Itoa(n)+"x  "+strconv.Itoa(sleepx), t, func() {

		Convey(" h1 fixed muti "+strconv.Itoa(n)+"x  "+strconv.Itoa(sleepx), func() {
			//maxFails := getAppMaxFails(conf, h1.AppName)

			isOpen := d.IsOpenFailedSleep(h1)
			So(isOpen, ShouldBeFalse)
			//
			isStop := false
			//模拟每秒有3~4个错误
			go func() {
				for {
					if isStop {
						break
					}
					m.Mark(1)
					time.Sleep(300 * time.Millisecond)
				}
			}()

			time.Sleep(sleep)

			fmt.Println(m.Snapshot().Rate1x())
			isOpen = d.IsOpenFailedSleep(h1)
			So(isOpen, ShouldBeTrue)
			So(h1.IsFailedSleepOpen, ShouldBeTrue)
			So(h1.LastSleepOpenTime, ShouldNotBeNil)
			So(h1.LastSleepExpectedCloseTime, ShouldNotBeNil)
			interval := h1.LastSleepExpectedCloseTime.UnixNano() - h1.LastSleepOpenTime.UnixNano()
			So(interval*x, ShouldBeGreaterThanOrEqualTo, seconds)

			isOpen = d.IsOpenFailedSleep(h1)
			So(isOpen, ShouldBeTrue)
			So(h1.IsFailedSleepOpen, ShouldBeTrue)
			So(h1.LastSleepOpenTime, ShouldNotBeNil)
			So(h1.LastSleepExpectedCloseTime, ShouldNotBeNil)
			interval = h1.LastSleepExpectedCloseTime.UnixNano() - h1.LastSleepOpenTime.UnixNano()
			So(interval*x, ShouldBeGreaterThanOrEqualTo, seconds)
			isStop = true
			time.Sleep(time.Duration(x) * sleep)
			isOpen = d.IsOpenFailedSleep(h1)
			So(isOpen, ShouldBeFalse)
			So(h1.IsFailedSleepOpen, ShouldBeFalse)
			So(h1.LastSleepOpenTime, ShouldBeNil)
			So(h1.LastSleepExpectedCloseTime, ShouldBeNil)

			//
			isStop = false
			//模拟每秒有1~2个错误
			go func() {
				for {
					if isStop {
						break
					}
					m.Mark(1)
					time.Sleep(900 * time.Millisecond)
				}
			}()
			isOpen = d.IsOpenFailedSleep(h1)
			So(isOpen, ShouldBeFalse)
			So(h1.IsFailedSleepOpen, ShouldBeFalse)
			So(h1.LastSleepOpenTime, ShouldBeNil)
			So(h1.LastSleepExpectedCloseTime, ShouldBeNil)

			isStop = true
			//wait for stop
			time.Sleep(500 * time.Millisecond)
			//模拟每秒有3个以上错误
			isStop = false
			go func() {
				for {
					if isStop {
						break
					}
					m.Mark(1)
					time.Sleep(200 * time.Millisecond)
				}
			}()
			time.Sleep(time.Duration(x) * sleep)
			isOpen = d.IsOpenFailedSleep(h1)
			So(isOpen, ShouldBeTrue)
			So(h1.IsFailedSleepOpen, ShouldBeTrue)
			So(h1.LastSleepOpenTime, ShouldNotBeNil)
			So(h1.LastSleepExpectedCloseTime, ShouldNotBeNil)
			interval = h1.LastSleepExpectedCloseTime.UnixNano() - h1.LastSleepOpenTime.UnixNano()
			So(interval*x, ShouldBeGreaterThanOrEqualTo, seconds)

			isStop = true

		})

	})
}

func TestDiscoveryBalancer_IsOpenFailedSleep_ForSeqMode(t *testing.T) {

	conf := kvs.NewEmptyMapConfigSource("map")
	conf.Set("lb.default.balancer.name", "WeightRobinRound")
	conf.Set("lb.default.max.fails", "3")
	conf.Set("lb.default.fail.time.window", "1s")
	conf.Set("lb.default.fail.sleep.mode", "seq")
	conf.Set("lb.default.fail.sleep.x", "1,1,2,3,5,8,13,21")
	conf.Set("lb.default.fail.sleep.max", "6s")
	d := NewDiscoveryBalancer(conf)
	//
	appName := "testForOpenFailedSleep"

	address, port := "192.168.1.2", "8081"
	h1 := newHostInstance(appName, address, port)
	d.AddHostInstance(appName, h1)
	//
	seconds := getAppFailTimeWindowSeconds(conf, h1.AppName)
	m := GetOrRegisterErrorMeter(h1, seconds)
	xs := getAppFailSleepSeqX(conf, h1.AppName)
	sleep := getAppMaxFailedSleep(conf, h1.AppName)
	//
	Convey("test IsOpenFailedSleep Seq Mode", t, func() {

		Convey(" h1 seq mode", func() {
			//maxFails := getAppMaxFails(conf, h1.AppName)

			isOpen := d.IsOpenFailedSleep(h1)
			So(isOpen, ShouldBeFalse)
			//
			isStop := false
			//模拟每秒有3~4个错误
			go func() {
				for {
					if isStop {
						break
					}
					m.Mark(1)
					time.Sleep(300 * time.Millisecond)
				}
			}()
			for i := 0; i < 10; i++ {
				ct := 0
				x := int64(xs[ct])
				sleepSeconds := int64(sleep.Seconds()) * x
				time.Sleep(time.Duration(sleepSeconds) * time.Second)
				fmt.Println(m.Snapshot().Rate1x())
				isOpen = d.IsOpenFailedSleep(h1)
				So(isOpen, ShouldBeTrue)
				So(h1.IsFailedSleepOpen, ShouldBeTrue)
				So(h1.LastSleepOpenTime, ShouldNotBeNil)
				So(h1.LastSleepExpectedCloseTime, ShouldNotBeNil)
				interval := h1.LastSleepExpectedCloseTime.UnixNano() - h1.LastSleepOpenTime.UnixNano()
				So(interval*x, ShouldBeGreaterThanOrEqualTo, seconds)
				ct++
			}

			isStop = true

		})

	})
}
