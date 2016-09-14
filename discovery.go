package zuul

import (
    "encoding/xml"
    "fmt"
    "io/ioutil"
    "net/http"
    "time"

    "github.com/tietang/go-eureka-client/eureka"
)

type Discovery struct {
    apps      *eureka.Applications
    AppNames  map[string]string
    eurekaUrl string
}

func NewDiscovery(eurekaUrl string) *Discovery {
    return &Discovery{eurekaUrl: eurekaUrl}
}

func (d *Discovery) ScheduleAtFixedRate(second time.Duration) {
    d.run()
    go d.runTask(second)
}

func (d *Discovery) runTask(second time.Duration) {
    timer := time.NewTicker(second)
    for {
        select {
        case <-timer.C:
            go d.run()
        }
    }
}

func (d *Discovery) run() {
    apps, err := d.GetApplications()
    if err == nil {
        d.apps = apps
    } else {
        fmt.Println(err)
    }
}

func (c *Discovery) GetApplications() (*eureka.Applications, error) {
    url := c.eurekaUrl + "/apps"

    //	req, err := http.NewRequest("GET", url, nil)
    //	req.Header.Add("Accept", "application/json")
    //	res, err := c.client.Do(req)
    //	http.Client.Do(req)
    res, err := http.Get(url)
    if err != nil {
        fmt.Println(err)
        return nil, err
    }
    //	fmt.Println(res.StatusCode)
    respBody, err := ioutil.ReadAll(res.Body)
    if err != nil {
        fmt.Println(err)
        return nil, err
    }
    if res.StatusCode != http.StatusOK {
        fmt.Println(err)
        return nil, err
    }
    var applications *eureka.Applications = new(eureka.Applications)
    err = xml.Unmarshal(respBody, applications)

    //	fmt.Println(string(respBody))
    //	fmt.Println(err, applications)
    return applications, err
}
