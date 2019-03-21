package discovery

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"k8s.io/client-go/rest"
	"testing"
)

func TestKubernetesDiscovery_GetService(t *testing.T) {
	config := &rest.Config{
		Host: "https://172.16.39.193:6443/",
	}
	kd := NewKubernetesDiscoveryByConfig(config)
	Convey("", t, func() {
		services, err := kd.GetServicesInTime()
		fmt.Println(err)
		for key, service := range services {
			fmt.Println(key, service)
		}

	})
}
