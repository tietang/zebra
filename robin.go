package zuul

type Balancer interface {
    Next(hosts []HostInstance) *HostInstance
}

type HostInstance struct {
    Id            int
    Name          string
    Weight        int
    CurrentWeight int
}

//片排序
type ByCurrentWeight []HostInstance

func (p ByCurrentWeight) Len() int {
    return len(p)
}
func (p ByCurrentWeight) Less(i, j int) bool {
    return p[i].CurrentWeight > p[j].CurrentWeight
}
func (p ByCurrentWeight) Swap(i, j int) {
    p[i], p[j] = p[j], p[i]
}

//type HostInstances []HostInstance

type Robin struct {
    TotalWeight int
    lastedHost  HostInstance
}

func NewRobinBalancer(hosts []HostInstance) *Robin {
    return &Robin{}
}

func (r *Robin) Next(hosts []HostInstance) *HostInstance {
    r.TotalWeight = r.totalWeight(hosts)
    totalWeight := r.TotalWeight
    for i := 0; i < len(hosts); i++ {
        hosts[i].CurrentWeight = hosts[i].CurrentWeight + hosts[i].Weight
    }
    //	sort.Sort(ByCurrentWeight(hosts))
    //	index := len(hosts) - 1
    //	hosts[index].CurrentWeight = (hosts[index].CurrentWeight - totalWeight)
    selected := r.max(hosts) //& hosts[index]
    selected.CurrentWeight = selected.CurrentWeight - totalWeight
    return selected
}

func (r *Robin) max(hosts []HostInstance) *HostInstance {
    max := &hosts[0]
    for i := 0; i < len(hosts); i++ {
        if hosts[i].CurrentWeight > max.CurrentWeight {
            max = &hosts[i]
        }
    }
    return max
}

func (r *Robin) totalWeight(hosts []HostInstance) int {
    totalWeight := 0
    for _, v := range hosts {
        totalWeight = totalWeight + v.Weight
    }
    return totalWeight
}
