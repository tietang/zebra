package router

import (
	"sync"
	"sync/atomic"
	"time"
)

type HostInstance struct {
	//负载均衡参数
	isInit        bool
	Weight        int32
	CurrentWeight int32
	EffectWeight  int32
	IsAvailable   bool
	//服务调用
	InstanceId           string
	Name                 string
	AppName              string
	AppGroupName         string
	DataCenter           string
	Tags                 []string
	Params               map[string]string
	ServiceSource        string
	ExternalInstance     interface{}
	Scheme               string
	Address              string
	Port                 string
	HealthCheckUrl       string
	Status               string
	LastUpdatedTimestamp int
	OverriddenStatus     string
	Version              string //实例应用版本号
	//
	LastSleepOpenTime          *time.Time //最后Sleep打开时间
	LastSleepExpectedCloseTime *time.Time //最后Sleep预计关闭时间
	LastSleepWindowSeqCount    int        //最后连续Sleep统计
	IsFailedSleepOpen          bool
	IsNotFound                 bool
	//
	lock *sync.Mutex
}

func (h *HostInstance) Init() {
	if h.lock == nil {
		h.lock = new(sync.Mutex)
	}

	if !h.isInit {
		h.EffectWeight = h.Weight
		h.CurrentWeight = h.Weight
		h.isInit = true
	}
}

func (h *HostInstance) ToKey() string {
	return h.InstanceId
}

func (h *HostInstance) SetCurrentWeight(newCurrentWeight int32) {
	//h.lock.Lock()
	//defer h.lock.Unlock()
	//h.CurrentWeight = CurrentWeight
	//atomic.StoreInt32(&h.CurrentWeight, newCurrentWeight)
	atomic.SwapInt32(&h.CurrentWeight, newCurrentWeight)
}

func (h *HostInstance) ResetFailedSleepFlag() {
	h.lock.Lock()
	defer h.lock.Unlock()
	h.LastSleepOpenTime = nil
	h.LastSleepExpectedCloseTime = nil
	h.LastSleepWindowSeqCount = 0
	h.IsFailedSleepOpen = false
}

func (h *HostInstance) SetFailedSleepFlagForSeqMode(sleepDuration time.Duration) {
	h.lock.Lock()
	defer h.lock.Unlock()
	lastSleepOpenTime := time.Now()
	//如果多次连续Sleep的第一次，则设置第一次Sleep时间
	if h.LastSleepOpenTime == nil {
		h.LastSleepOpenTime = &lastSleepOpenTime
	} else {
		lastSleepOpenTime = *h.LastSleepOpenTime
	}
	closeTime := lastSleepOpenTime.Add(sleepDuration)
	h.LastSleepExpectedCloseTime = &closeTime
	h.LastSleepWindowSeqCount++
	h.IsFailedSleepOpen = true
}

func (h *HostInstance) SetFailedSleepFlagForFixedMode(sleepDuration time.Duration) {
	h.lock.Lock()
	defer h.lock.Unlock()
	lastSleepOpenTime := time.Now()
	h.LastSleepOpenTime = &lastSleepOpenTime
	closeTime := lastSleepOpenTime.Add(sleepDuration)
	h.LastSleepExpectedCloseTime = &closeTime

	h.LastSleepWindowSeqCount++
	h.IsFailedSleepOpen = true
}

//复制原生属性
func (h *HostInstance) CopyFrom(source *HostInstance) {
	h.InstanceId = source.InstanceId
	h.Name = source.Name
	h.AppName = source.AppName
	h.AppGroupName = source.AppGroupName
	h.DataCenter = source.DataCenter
	//h.isInit = source.isInit
	if h.EffectWeight == h.Weight {
		h.EffectWeight = source.Weight
	}
	h.Weight = source.Weight
	h.Tags = source.Tags
	h.Params = source.Params
	h.ServiceSource = source.ServiceSource
	h.ExternalInstance = source.ExternalInstance
	h.Scheme = source.Scheme
	h.Address = source.Address
	h.Port = source.Port
	h.HealthCheckUrl = source.HealthCheckUrl
	h.Status = source.Status
	h.LastUpdatedTimestamp = source.LastUpdatedTimestamp
	h.OverriddenStatus = source.OverriddenStatus
}

//片排序
type ByCurrentWeight []*HostInstance

func (p ByCurrentWeight) Len() int {
	return len(p)
}

func (p ByCurrentWeight) Less(i, j int) bool {
	return p[i].CurrentWeight >= p[j].CurrentWeight
}

func (p ByCurrentWeight) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
