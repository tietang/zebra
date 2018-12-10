package main
//
//import (
//    "fmt"
//    "github.com/tietang/zebra/discovery"
//    "k8s.io/client-go/rest"
//)
//
//func main() {
//    config := &rest.Config{
//        Host: "http://172.16.30.112:8080/",
//    }
//    kd := discovery.NewKubernetesDiscoveryByConfig(config)
//    services, err := kd.GetServicesInTime()
//    fmt.Println(err)
//    for key, service := range services {
//        fmt.Println(key, "")
//        for _, ins := range service.Instances {
//            fmt.Println("", ins)
//        }
//    }
//
//}
