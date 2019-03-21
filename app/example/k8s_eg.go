package main

import (
	"encoding/json"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	//http://172.16.30.112:8080/api/v1/namespaces/default/services/nginx-service
	//http://172.16.30.112:8080/api/v1/namespaces/default/endpoints/nginx-service
	config := &rest.Config{
		Host: "http://172.16.30.112:8080/",
	}
	c, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	list, err := c.CoreV1().Services("").List(metav1.ListOptions{})
	fmt.Println(err)
	fmt.Println(list)
	for _, s := range list.Items {
		name := s.ObjectMeta.Name
		namespace := s.ObjectMeta.Namespace
		endpoints, err := c.CoreV1().Endpoints(namespace).Get(name, metav1.GetOptions{})
		if err != nil {
			continue
		}
		d, _ := json.Marshal(s)
		fmt.Println(string(d))
		for _, set := range endpoints.Subsets {
			for _, address := range set.Addresses {
				d, _ := json.Marshal(address)
				fmt.Println("   " + string(d))
			}

		}
	}

	//endpointsList, err := c.CoreV1().Endpoints("").List(metav1.ListOptions{})
	//if err != nil {
	//    //return nil, err
	//}
	//
	//for _, endpoints := range endpointsList.Items {
	//    name := endpoints.Name
	//    q := metav1.ListOptions{
	//        FieldSelector: "metadata.name=" + name,
	//    }
	//    list, err := c.CoreV1().Services(endpoints.Namespace).List(q)
	//    fmt.Println(err)
	//    d, _ := json.Marshal(list.Items)
	//    fmt.Println(string(d))
	//    for _, s := range endpoints.Subsets {
	//        d, _ := json.Marshal(s)
	//        fmt.Println(string(d))
	//    }
	//}
}
