package main

import "net/http"

func main() {
    //	eurekaUrl := "http://dev.discovery.shishike.com/eureka"
    //	eurekaUrl := "http://discovery2.keruyun.com:8761/eureka"
    //	robin := zuul.NewDiscoveryRobin(eurekaUrl, &zuul.Robin{})
    //	robin.ScheduleAtFixedRate(10 * time.Second)
    //	//	apps, err := robin.GetApplications()

    //	for i := 0; i <= 10; i++ {
    //		app, ins := robin.Next("DISCOVERY")
    //		//		fmt.Println(app, in/s)
    //		fmt.Println(app.Name, ins)
    //		//		fmt.Println(app.Instances[ins.Id])
    //	}

    var f func()
    f = func() {}
    println(f == nil)

    http.ServeMux
}
