package discovery

import (
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"strconv"
	"strings"
	"sync"
	"time"
)

type KubernetesDiscovery struct {
	services  map[string]*Service
	callbacks []func(map[string]*Service)
	Config    *rest.Config
	Client    *kubernetes.Clientset
	lock      *sync.Mutex
}

func NewKubernetesDiscoveryByConfig(config *rest.Config) *KubernetesDiscovery {

	c := config
	if config == nil {
		cf, err := rest.InClusterConfig()
		if err != nil {
			panic(err)
		}
		c = cf
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(c)
	if err != nil {
		panic(err)
	}
	//list :=
	//    clientset.CoreV1().Endpoints("").List(metav1.ListOptions{})
	cd := &KubernetesDiscovery{Config: config, callbacks: make([]func(map[string]*Service), 0)}
	cd.Client = clientset
	cd.lock = new(sync.Mutex)
	return cd
}

//address: [http[s]://hostName:port]
func NewKubernetesDiscovery(address string) *KubernetesDiscovery {
	config := &rest.Config{
		Host: address,
	}
	//config, err := rest.InClusterConfig()
	//if err != nil {
	//    panic(err)
	//}
	//
	//config.Host = address
	return NewKubernetesDiscoveryByConfig(config)
}

func (d *KubernetesDiscovery) AddCallback(callback func(map[string]*Service)) {
	d.callbacks = append(d.callbacks, callback)

}
func (c *KubernetesDiscovery) GetServicesInTime() (map[string]*Service, error) {
	serviceList, err := c.Client.CoreV1().Services("").List(metav1.ListOptions{})

	if err != nil {
		return nil, err
	}
	services := make(map[string]*Service, 0)

	for _, s := range serviceList.Items {
		name := s.ObjectMeta.Name
		labels := s.Labels
		if labels == nil {
			labels = make(map[string]string)
		}
		namespace := s.ObjectMeta.Namespace
		service := &Service{
			Name:      name,
			Labels:    labels, //prefix
			Instances: make([]*Instance, 0),
		}
		port := s.Spec.Ports[0].NodePort
		if port <= 0 {
			continue
			//port = s.Spec.Ports[0].Port
		}
		//if port <= 0 {
		//    continue
		//}
		//
		endpoints, err := c.Client.CoreV1().Endpoints(namespace).Get(name, metav1.GetOptions{})
		if err != nil {
			continue
		}
		if endpoints.Namespace != namespace {
			continue
		}
		//d, _ := json.Marshal(s)
		//fmt.Println(string(d))
		for _, set := range endpoints.Subsets {
			for _, address := range set.Addresses {

				if port <= 0 {
					port = set.Ports[0].Port
				}
				portStr := strconv.Itoa(int(port))
				addr := address.NodeName
				if addr == nil {
					addr = &address.IP
				}

				ins := &Instance{
					Name:       name,
					AppName:    name,
					InstanceId: strings.Join([]string{name, *addr, portStr}, ":"),
					Port:       portStr,
					Address:    *addr,
					Params:     s.Labels,
				}
				c.addInstance(service, ins)
				//fmt.Println("   ", endpoints, ins)
			}
		}
		if len(service.Instances) > 0 {
			services[name] = service
		}

	}

	return services, nil

}

func (c *KubernetesDiscovery) addInstance(service *Service, ins *Instance) {
	c.lock.Lock()
	defer c.lock.Unlock()
	isExists := false
	for i, instance := range service.Instances {
		if ins.InstanceId == instance.InstanceId {
			service.Instances[i] = ins
			isExists = true
		}
	}
	if !isExists {
		service.Instances = append(service.Instances, ins)
	}
}

func (c *KubernetesDiscovery) add(services map[string]*Service, name string, ins *Instance) {
	c.lock.Lock()
	defer c.lock.Unlock()
	service, ok := services[name]
	isExists := false
	if ok {
		for i, instance := range service.Instances {
			if ins.InstanceId == instance.InstanceId {
				service.Instances[i] = ins
				isExists = true
			}
		}
		services[name] = service
	} else {
		service = &Service{
			Name:      name,
			Instances: make([]*Instance, 0),
		}

	}
	if !isExists {
		service.Instances = append(service.Instances, ins)
		services[name] = service
	}
}

func (c *KubernetesDiscovery) GetServices() map[string]*Service {
	if c.services == nil {
		services, err := c.GetServicesInTime()
		if err == nil {
			return services
		}
	}
	return c.services
}

func (c *KubernetesDiscovery) GetService(name string) *Service {
	if c.services == nil {
		log.Info("catalogServices is nil")
		return nil
	}
	for name, service := range c.services {
		if strings.ToLower(name) == strings.ToLower(name) {
			return service
		}
	}
	return nil
}

func (d *KubernetesDiscovery) Watching(second time.Duration) {
	d.run()
	go d.runTask(second)
}

func (d *KubernetesDiscovery) runTask(second time.Duration) {
	timer := time.NewTicker(second)
	for {
		select {
		case <-timer.C:
			go d.run()
		}
	}
}
func (d *KubernetesDiscovery) run() {
	services, err := d.GetServicesInTime()
	if err == nil || services != nil {
		d.services = services
		d.execCallbacks(d.services)
	} else {
		log.Error(err)
	}
}

func (d *KubernetesDiscovery) execCallbacks(services map[string]*Service) {
	if len(d.callbacks) > 0 {
		for _, c := range d.callbacks {
			go c(services)
		}
	}
}

func (c *KubernetesDiscovery) Health() (bool, string) {
	c.Client.CoreV1().RESTClient().APIVersion()
	//if err != nil || leader == "" {
	//    return false, err.Error()
	//}

	ok, desc := true, "ok"

	return ok, desc

}
