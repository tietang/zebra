package router

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"encoding/json"
	"github.com/rcrowley/go-metrics"
	log "github.com/sirupsen/logrus"
	"github.com/tietang/props/kvs"
	"github.com/tietang/zebra/utils"
)

const (
	STATUS_UP   = "UP"   //Ready to receive traffic
	STATUS_DOWN = "DOWN" // Do not send traffic- healthcheck routeChangedCallback failed

	//
	LB_LOG_ENABLED = "lb.log.selected.enabled"
	LB_LOG_ALL     = "lb.log.selected.all"

	LB_DEFAULT_NAME             = "lb.default"
	DEFAULT_LB_MAX_FAILS        = 3
	DEFAULT_FAIL_TIME_WINDOW    = 10 * time.Second
	DEFAULT_FAIL_SLEEP_X        = 1
	DEFAULT_FAIL_SLEEP_MODE_SEQ = "seq"
	DEFAULT_FAIL_SLEEP_MODE     = DEFAULT_FAIL_SLEEP_MODE_SEQ
	DEFAULT_FAIL_SLEEP_MAX      = 60 * time.Second
	DEFAULT_LB_NAME             = "WeightRobinRound"
	INSTANCE_EXPIRED_TIME       = 60 * time.Second
	DEFAULT_FIBONACCI_BASE      = 100
)

var DEFAULT_FAIL_SLEEP_SEQ_X = []int{1, 1, 2, 3, 5, 8}

type DiscoveryBalancer struct {
	Balancer            Balancer
	Hosts               *sync.Map // map[string][]*HostInstance
	UnavailableHosts    *sync.Map // 只用于主动监控检查 map[string][]*HostInstance
	HostInstanceSources []HostInstanceSource
	lock                *sync.RWMutex
	conf                kvs.ConfigSource
	ServiceBalancers    map[string]Balancer
	//
	SleepHostInstances []*HostInstance
	healthCheckTicker  *time.Ticker
}

func NewDiscoveryBalancer(conf kvs.ConfigSource) *DiscoveryBalancer {
	d := &DiscoveryBalancer{
		Hosts:              new(sync.Map), //make(map[string][]*HostInstance),
		UnavailableHosts:   new(sync.Map), //make(map[string][]*HostInstance),
		lock:               new(sync.RWMutex),
		conf:               conf,
		ServiceBalancers:   make(map[string]Balancer),
		SleepHostInstances: make([]*HostInstance, 0),
	}
	d.initDefaultBalancer()

	return d
}

func (r *DiscoveryBalancer) initDefaultBalancer() {
	r.Balancer = r.GetBalancer(LB_DEFAULT_NAME)
	//开启健康检测ticker
	r.healthCheckTicker = time.NewTicker(1 * time.Second)
	go func() {
		for {
			select {
			case <-r.healthCheckTicker.C:
				r.UnavailableHosts.Range(func(key, value interface{}) bool {
					r.healthCheck(key.(string), value.([]*HostInstance))
					return true
				})
			}
		}
	}()

}

func (r *DiscoveryBalancer) healthCheck(key string, hs []*HostInstance) bool {

	cleanHosts := make([]*HostInstance, 0)
	for _, h := range hs {

		if h.HealthCheckUrl != "" {
			res, err := http.Get(h.HealthCheckUrl)
			//健康检查通过，使得实例状态为可用
			if err == nil && res.StatusCode == http.StatusOK {
				h.IsAvailable = true
			} else {
				cleanHosts = append(cleanHosts, h)
			}

			log.Info("check health: ", res, h.HealthCheckUrl)
		}
	}
	r.UnavailableHosts.Store(key, cleanHosts)
	return false
}

func (r *DiscoveryBalancer) GetBalancer(serviceId string) Balancer {
	if balancer, ok := r.ServiceBalancers[serviceId]; ok {
		return balancer
	}
	//name := strings.ToLower(r.conf.GetDefault(fmt.Sprintf(KEY_LB_NAME_TEMPLATE, serviceId), DEFAULT_LB_NAME))
	defName := r.conf.GetDefault(fmt.Sprintf(KEY_LB_NAME_TEMPLATE, LB_DEFAULT_NAME), DEFAULT_LB_NAME)
	name := r.conf.GetDefault(fmt.Sprintf(KEY_LB_NAME_TEMPLATE, serviceId), defName)
	defBaseVal := r.conf.GetIntDefault(fmt.Sprintf(KEY_FIBONACCI_BASE_TEMPLATE, LB_DEFAULT_NAME), DEFAULT_FIBONACCI_BASE)
	base := r.conf.GetIntDefault(fmt.Sprintf(KEY_FIBONACCI_BASE_TEMPLATE, serviceId), defBaseVal)
	log.Info("DEBUG service id: ", serviceId, name)
	balancer := DefaultBalancerFactory.New(name)
	if name == "fibonacci" {
		fibonacci := balancer.(*FibonacciWeightRobinRound)
		fibonacci.Base = int64(base)
	}
	r.ServiceBalancers[serviceId] = balancer
	return balancer
}

func (r *DiscoveryBalancer) register(hostInstanceSource HostInstanceSource) {
	if hostInstanceSource == nil {
		return
	}
	hostInstanceSource.SetAppHostsChangedCallback(func(appName string, hs []*HostInstance) {
		r.AddOrUpdateHostInstances(appName, hs)
	})
	r.HostInstanceSources = append(r.HostInstanceSources, hostInstanceSource)
	log.Info("HostInstanceSource register: ", hostInstanceSource.Name())
}

//
func (d *DiscoveryBalancer) Start() {
	for _, ins := range d.HostInstanceSources {
		log.Info("HostInstanceSource: ", ins.Name(), " is starting...")
		ins.Start()
	}
}

func (d *DiscoveryBalancer) AddOrUpdateHostInstances(appName string, hs []*HostInstance) {
	//debugHosts(d.Hosts, "0: ")
	//debugHosts(d.UnavailableHosts, "0: ")
	if len(hs) == 0 {
		return
	}
	//初始化balancer
	d.GetBalancer(appName)
	key := strings.ToUpper(appName)
	for _, his := range hs {
		d.add(d.Hosts, appName, his)
	}
	val, ok := d.Hosts.Load(key)
	var hosts []*HostInstance
	if ok {
		hosts = val.([]*HostInstance)
	} else {
		hosts = make([]*HostInstance, 0)
	}

	cleanHosts := make([]*HostInstance, 0)
	//删除过期实例
	for _, his := range hosts {
		//如果很久没被LB命中，或者在当前可用列表存在但不在服务发现服务器端存在，则被删除
		if !d.contains(hs, his) {
			continue
		}
		cleanHosts = append(cleanHosts, his)
	}

	d.Hosts.Store(key, cleanHosts)
	//可用列表
	valHosts, ok := d.Hosts.Load(key)
	var aHosts []*HostInstance
	if ok {
		aHosts = valHosts.([]*HostInstance)
	} else {
		aHosts = make([]*HostInstance, 0)
	}

	//debugHosts(d.Hosts, "3: ")
	//debugHosts(d.UnavailableHosts, "3: ")

	//更新或过滤不可用列表
	d.updateUnavailableHostInstance(appName, aHosts)

	//debugHosts(d.Hosts, "4: ")
	//debugHosts(d.UnavailableHosts, "4: ")
}

//
//func debugHosts(m *sync.Map, i string) {
//    lm := make(map[string]interface{}, )
//    size := 0
//    lens := make(map[string]int)
//    m.Range(func(key, value interface{}) bool {
//        lm[key.(string)] = value
//        //xhosts := value.([]*HostInstance)
//        size++
//        if l, ok := lens[key.(string)]; ok {
//            l++
//            lens[key.(string)] = l
//        } else {
//            lens[key.(string)] = 1
//        }
//        return true
//    })
//
//    d1, _ := json.Marshal(lm)
//
//    log.Debug(i, size, " ", lens, "   ", string(d1))
//
//}

func (d *DiscoveryBalancer) contains(hs []*HostInstance, h *HostInstance) bool {
	for _, his := range hs {
		if his.InstanceId == h.InstanceId {
			return true
		}
	}
	return false
}

func (d *DiscoveryBalancer) updateUnavailableHostInstance(appName string, hs []*HostInstance) {
	key := strings.ToUpper(appName)
	d.lock.Lock()
	defer d.lock.Unlock()

	val, ok := d.UnavailableHosts.Load(key)
	var ins []*HostInstance
	if ok {
		ins = val.([]*HostInstance)
	} else {
		ins = make([]*HostInstance, 0)
	}

	newHosts := make([]*HostInstance, 0)
	//如果存在就替换
	for _, host := range ins {
		for _, hi := range hs {
			if host.InstanceId == hi.InstanceId {
				newHosts = append(newHosts, host)
			}
		}
	}
	d.UnavailableHosts.Store(key, newHosts)

}

func (d *DiscoveryBalancer) removeInstanceFromUnavailable(appName string, h *HostInstance) {
	key := strings.ToUpper(appName)

	val, ok := d.UnavailableHosts.Load(key)
	var ins []*HostInstance
	if ok {
		ins = val.([]*HostInstance)
	} else {
		ins = make([]*HostInstance, 0)
	}
	cleanHosts := make([]*HostInstance, 0)
	for _, hi := range ins {
		if hi.InstanceId != h.InstanceId {
			cleanHosts = append(cleanHosts, hi)
		}
	}

	d.UnavailableHosts.Store(key, cleanHosts)
}

func (d *DiscoveryBalancer) AddHostInstance(appName string, h *HostInstance) *HostInstance {
	return d.add(d.Hosts, appName, h)
}

func (d *DiscoveryBalancer) AddUnavailableHostInstance(appName string, h *HostInstance) {

	d.add(d.UnavailableHosts, appName, h)
}

func (d *DiscoveryBalancer) add(hosts *sync.Map, appName string, h *HostInstance) *HostInstance {
	key := strings.ToUpper(appName)
	d.lock.Lock()
	defer d.lock.Unlock()

	value, ok := hosts.Load(key)
	var ins []*HostInstance
	if ok {
		ins = value.([]*HostInstance)
	} else {
		ins = make([]*HostInstance, 0)
	}
	hasExists := false

	//如果存在就替换
	for i, host := range ins {
		if host.InstanceId == h.InstanceId {
			hasExists = true
			isAvailable := host.IsAvailable
			ins[i] = h //.CopyFrom(h)
			ins[i].ResetFailedSleepFlag()
			d.removeInstanceFromUnavailable(appName, h)
			log.Debug("updated host instance: ", appName, ", ", h.InstanceId)
			ins[i].IsAvailable = isAvailable
		}
	}
	//debugHosts(d.Hosts, "1: ")
	//debugHosts(d.UnavailableHosts, "1: ")
	//不存在，就追加
	if !hasExists {
		h.IsAvailable = true
		ins = append(ins, h)
		log.Debug("add host instance: ", appName, ", ", h.InstanceId)
	}

	hosts.Store(key, ins)
	//debugHosts(d.Hosts, "2: ")
	//debugHosts(d.UnavailableHosts, "2: ")
	return h

}

func (d *DiscoveryBalancer) Next(appName string, key string, matched bool) *HostInstance {
	h := d.next(appName, key, matched)
	//if d.conf.GetBoolDefault(LB_LOG_ENABLED, false) && h != nil {
	//    if d.conf.GetBoolDefault(LB_LOG_ALL, false) {
	//        log.WithField("selected", h).Info()
	//    } else {
	//        log.Info("selected: ", h.InstanceId, ", ", h.Address)
	//    }
	//}
	//log.WithField("ins", h).Info("DEBUG")

	//appKey := strings.ToUpper(appName)
	//value, ok := d.Hosts.Load(appKey)
	//var ins []*HostInstance
	//if ok {
	//    ins = value.([]*HostInstance)
	//} else {
	//    ins = make([]*HostInstance, 0)
	//}
	//log.Info(len(ins))
	//for _, h := range ins {
	//    log.WithField("ins", h).Info()
	//}

	return h
}

//返回nil表示无可用instance
func (d *DiscoveryBalancer) next(appName, key string, matched bool) *HostInstance {
	appKey := strings.ToUpper(appName)

	value, ok := d.Hosts.Load(appKey)

	if !ok {
		return nil
	}
	hosts := value.([]*HostInstance)
	isEnabledGray := d.conf.GetBoolDefault(fmt.Sprintf("traffic.cond.%s.enabled", appName), false)
	if isEnabledGray {
		grayhosts := make([]*HostInstance, 0)
		ahosts := make([]*HostInstance, 0)
		version := d.conf.GetDefault(fmt.Sprintf("traffic.cond.%s.version", appName), "non-version")
		for _, host := range hosts {
			if host.Version == version {
				grayhosts = append(grayhosts, host)
			} else {
				ahosts = append(ahosts, host)
			}
		}
		//流量切分实例过滤
		if matched {
			hosts = grayhosts
		} else {
			hosts = ahosts
		}
	}

	val, ok := d.UnavailableHosts.Load(appKey)
	var uahosts []*HostInstance
	if ok {
		uahosts = val.([]*HostInstance)
	} else {
		uahosts = make([]*HostInstance, 0)
	}

	if len(hosts) == 0 {
		return nil
	}
	//if !ok && hi.Status == "UP" {
	//    return hi
	//}
	balancer := d.GetBalancer(appName)
	h := d.nextHostInstance(balancer, key, hosts, uahosts, nil)
	return h
}

func (r *DiscoveryBalancer) getAppFailSleepSeqX(name string) []int {

	return getAppFailSleepSeqX(r.conf, name)
}
func (r *DiscoveryBalancer) getAppFailSleepX(name string) int {
	return getAppFailSleepX(r.conf, name)
}

func (r *DiscoveryBalancer) getAppFailSleepMode(name string) string {
	return getAppFailSleepMode(r.conf, name)
}
func (r *DiscoveryBalancer) getAppFailTimeWindow(name string) time.Duration {
	return getAppFailTimeWindow(r.conf, name)
}

func (r *DiscoveryBalancer) getResponseTime(host string) (int64, bool) {
	timer := metrics.GetOrRegisterTimer(host, InstanceRegistry).Snapshot()
	if timer == nil {
		return 1, false
	}
	time := int64(timer.Mean())
	if time <= 0 {
		return 1, false
	}
	return time, true
}

type NextStack struct {
	filteredHosts *utils.Set
	depth         int
}

//TODO bug 当有某个实例不在线时
//filteredHosts 当所有实例不可用，并且还未加入到不可用实例，防止递归死循环

func (d *DiscoveryBalancer) nextHostInstance(balancer Balancer, key string, hosts []*HostInstance, unavailableHosts []*HostInstance, nextStack *NextStack) *HostInstance {
	//
	if nextStack == nil {
		nextStack = &NextStack{
			filteredHosts: utils.NewSet(),
			depth:         0,
		}
	}
	//如果可用数和不可用数一致，则认为无可用实例，直接返回nil
	if len(hosts) <= len(unavailableHosts) || len(hosts) <= nextStack.filteredHosts.Size() {
		return nil
	}
	//防止递归bug致死
	if nextStack.depth >= 100 {
		log.Warn("depth>=", nextStack.depth, "，may be a bug.")
		return nil
	}
	//按照配置的负载均衡算法，获取一个实例
	hi := balancer.Next(key, hosts)
	if hi == nil {
		return hi
	}
	nextStack.depth++
	//如果实例被Sleep或者不在线，则重新选择一个可用实例
	if d.IsOpenFailedSleep(hi) || hi.Status != STATUS_UP || !hi.IsAvailable {
		data, _ := json.Marshal(hi)
		log.Debug(string(data))
		nextStack.filteredHosts.Add(hi)
		//加入不可用列表
		d.AddUnavailableHostInstance(hi.AppName, hi)
		hi.IsAvailable = false
		return d.nextHostInstance(balancer, key, hosts, unavailableHosts, nextStack)
	}

	return hi
}

//KEY_FAIL_TIME_WINDOW_TEMPLATE = "%s.fail.time.window"
//KEY_FAIL_SLEEP_MODE_TEMPLATE = "%s.fail.sleep.mode"
//KEY_FAIL_SLEEP_X_TEMPLATE = "%s.fail.sleep.x"
//KEY_FAIL_SLEEP_MAX_TEMPLATE = "%s.fail.sleep.max"

//Sleep监测
func (r *DiscoveryBalancer) IsOpenFailedSleep(ins *HostInstance) bool {
	//仍然在Sleep窗口
	now := time.Now()
	if ins.LastSleepExpectedCloseTime != nil && ins.LastSleepExpectedCloseTime.After(now) {
		log.WithFields(log.Fields{
			"HostInstance": ins,
		}).Error("to be sleep cause for still in Sleep time window： ")
		return true
	}
	//超过Sleep关闭时间，则检测失败
	//检测当前app的当前时间窗口错误次数
	mode := r.getAppFailSleepMode(ins.AppName)
	if strings.ToLower(mode) == DEFAULT_FAIL_SLEEP_MODE_SEQ {
		return r.isFailedSleepForSeqMode(ins)
	} else {
		return r.isFailedSleepForFixedMode(ins)
	}
}

func (r *DiscoveryBalancer) isFailedSleepForFixedMode(ins *HostInstance) bool {

	//配置的Sleep时间
	d := getAppFailTimeWindowSeconds(r.conf, ins.AppName)
	//当前错误数
	m := GetOrRegisterErrorMeterSnapshot(ins, d)
	errCount := int(m.Rate1x() + 0.5)
	//default.max.fails=lb.default.max.fails
	//配置的时间窗口内最大失败次数
	fails := getAppMaxFails(r.conf, ins.AppName)

	//log.Info(errCount, " ", fails)
	//错误数大于配置阈值，Sleep
	if errCount >= fails {
		log.WithFields(log.Fields{
			"fails":        fails,
			"errCount":     errCount,
			"HostInstance": ins,
		}).Error("to be sleep cause for error greater than or equals to fails.")
		//配置的Sleep时间
		w := r.getAppFailTimeWindow(ins.AppName)
		//失败后，不服务sleep周期倍数
		x := r.getAppFailSleepX(ins.AppName)
		sleepDuration := utils.DurationMuti(w, int64(x))
		//如果sleepDuration大于配置的fail.sleep.max,则用fail.sleep.max替代
		appMaxFailedSleep := getAppMaxFailedSleep(r.conf, ins.AppName)
		if sleepDuration.Nanoseconds() >= appMaxFailedSleep.Nanoseconds() {
			sleepDuration = appMaxFailedSleep
		}
		ins.SetFailedSleepFlagForFixedMode(sleepDuration)
		return true
	} else {
		//如果在sleep窗口外无失败sleep，则重置
		if ins.LastSleepOpenTime != nil {
			log.WithFields(log.Fields{
				"fails":        fails,
				"errCount":     errCount,
				"HostInstance": ins,
			}).Info("to be back to normal.")
		}
		ins.ResetFailedSleepFlag()
	}

	return false
}

func (r *DiscoveryBalancer) isFailedSleepForSeqMode(ins *HostInstance) bool {

	//配置的Sleep时间
	d := getAppFailTimeWindowSeconds(r.conf, ins.AppName)
	//配置的Sleep时间
	w := r.getAppFailTimeWindow(ins.AppName)
	//当前错误数
	m := GetOrRegisterErrorMeterSnapshot(ins, d)
	errCount := int(m.Rate1x())
	//default.max.fails=lb.default.max.fails
	//配置的时间窗口内最大失败次数
	fails := getAppMaxFails(r.conf, ins.AppName)

	//如果当前时间已经超过预期的sleep窗口关闭时间1.5个周期，则作为新的Seq Sleep window开始；否则仍然为seq Sleep window期间。
	expectedNextWindow := time.Now().Truncate(utils.DurationOneHalf(w))
	if ins.LastSleepExpectedCloseTime != nil && ins.LastSleepExpectedCloseTime.Before(expectedNextWindow) {
		//如果不Sleep，则重置
		ins.ResetFailedSleepFlag()
		if ins.LastSleepOpenTime != nil {
			log.WithFields(log.Fields{
				"fails":        fails,
				"errCount":     errCount,
				"HostInstance": ins,
			}).Info("to be back to normal.")
		}
	}

	//错误数大于配置阈值，Sleep
	if errCount >= fails {
		log.WithFields(log.Fields{
			"fails":        fails,
			"errCount":     errCount,
			"HostInstance": ins,
		}).Error("to be sleep cause for error greater than or equals to fails.")
		//如果多次连续Sleep的第一次，则设置第一次Sleep时间
		if ins.LastSleepOpenTime == nil {
			now := time.Now()
			ins.LastSleepOpenTime = &now
		}

		//失败后，服务停止周期倍数
		x := r.getAppFailSleepX(ins.AppName)

		sleepx := r.getAppFailSleepSeqX(ins.AppName)
		limit := ins.LastSleepWindowSeqCount

		if limit >= len(sleepx) {
			limit = len(sleepx) - 1
		}
		x = 1
		for i := 0; i < limit; i++ {
			if len(sleepx) <= i {
				x += sleepx[i]
			}
		}
		x = sleepx[ins.LastSleepWindowSeqCount]
		sleepDuration := utils.DurationMuti(w, int64(x))
		//如果sleepDuration大于配置的fail.sleep.max,则用fail.sleep.max替代
		appMaxFailedSleep := getAppMaxFailedSleep(r.conf, ins.AppName)
		if sleepDuration.Nanoseconds() >= appMaxFailedSleep.Nanoseconds() {
			sleepDuration = appMaxFailedSleep
		}

		ins.SetFailedSleepFlagForSeqMode(sleepDuration)

		return true
	} else {
		//如果不Sleep，则重置
		if ins.LastSleepOpenTime != nil {
			log.WithFields(log.Fields{
				"fails":        fails,
				"errCount":     errCount,
				"HostInstance": ins,
			}).Info("to be back to normal.")
		}
		ins.ResetFailedSleepFlag()
	}

	return false
}

func (d *DiscoveryBalancer) stop() {
}
