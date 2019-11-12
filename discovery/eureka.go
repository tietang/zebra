package discovery

import (
    "encoding/xml"
    log "github.com/sirupsen/logrus"
    "github.com/tietang/go-eureka-client/eureka"
    "github.com/tietang/go-utils/errs"
    "io/ioutil"
    "net/http"
    "net/url"
    "strings"
    "sync/atomic"
    "time"
)

type EurekaDiscovery struct {
    apps      *eureka.Applications
    AppNames  map[string]string
    eurekaUrl []string
    callbacks []func(*eureka.Applications)
    ct        uint64
}

func NewEurekaDiscovery(eurekaUrl []string) *EurekaDiscovery {
    return &EurekaDiscovery{eurekaUrl: eurekaUrl, callbacks: make([]func(*eureka.Applications), 0)}
}

func (d *EurekaDiscovery) AddCallback(callback func(*eureka.Applications)) {
    d.callbacks = append(d.callbacks, callback)

}

func (d *EurekaDiscovery) execCallbacks(apps *eureka.Applications) {
    if len(d.callbacks) > 0 {
        for _, c := range d.callbacks {
            go c(apps)
        }
    }
}

func (d *EurekaDiscovery) ScheduleAtFixedRate(second time.Duration) {
    d.run()
    go d.runTask(second)
}

func (d *EurekaDiscovery) runTask(second time.Duration) {
    timer := time.NewTicker(second)
    for {
        select {
        case <-timer.C:
            go d.run()
        }
    }
}

func (d *EurekaDiscovery) run() {
    apps, err := d.GetApplications()
    if err == nil || apps != nil {
        d.apps = apps
        d.execCallbacks(apps)
    } else {
        log.Error(err)
    }
}

func (c *EurekaDiscovery) GetApps() *eureka.Applications {
    if c.apps == nil {
        apps, err := c.GetApplications()
        if err == nil {
            return apps
        }
    }
    return c.apps
}

func (c *EurekaDiscovery) GetApp(name string) *eureka.Application {
    if c.apps == nil {
        log.Info("Applications is nil")
        return nil
    }
    for _, app := range c.apps.Applications {
        if strings.ToLower(app.Name) == strings.ToLower(name) {
            return &app
        }
    }
    return nil
}

func (c *EurekaDiscovery) GetInstances(name string) (*eureka.Application, error) {
    //url := c.eurekaUrl + "/apps"
    url := c.getEurekaServerUrl() + "/apps/" + name

    //	req, err := http.NewRequest("GET", url, nil)
    //	req.Header.Add("Accept", "application/json")
    //	res, err := c.client.Do(req)
    //	http.Client.Do(req)
    res, err := http.Get(url)
    if err != nil {
        log.Error(err)
        return nil, err
    }
    //	log.Info(res.StatusCode)
    respBody, err := ioutil.ReadAll(res.Body)
    if err != nil {
        log.Error(err)
        return nil, err
    }
    if res.StatusCode != http.StatusOK {
        log.Error(err)
        return nil, err
    }
    app := new(eureka.Application)
    err = xml.Unmarshal(respBody, app)

    //	log.Info(string(respBody))
    //	log.Info(err, applications)
    return app, err
}

func (c *EurekaDiscovery) GetApplications() (*eureka.Applications, error) {
    //url := c.eurekaUrl + "/apps"
    url := c.getEurekaServerUrl() + "/apps"

    //	req, err := http.NewRequest("GET", url, nil)
    //	req.Header.Add("Accept", "application/json")
    //	res, err := c.client.Do(req)
    //	http.Client.Do(req)
    res, err := http.Get(url)
    if err != nil {
        log.Error(err)
        return nil, err
    }
    //	log.Info(res.StatusCode)
    respBody, err := ioutil.ReadAll(res.Body)
    if err != nil {
        log.Error(err)
        return nil, err
    }
    if res.StatusCode != http.StatusOK {
        log.Error(err)
        return nil, err
    }
    var applications *eureka.Applications = new(eureka.Applications)
    err = xml.Unmarshal(respBody, applications)

    //	log.Info(string(respBody))
    //	log.Info(err, applications)
    return applications, err
}

func (c *EurekaDiscovery) getEurekaServerUrl() string {
    ct := atomic.AddUint64(&c.ct, 1)
    size := len(c.eurekaUrl)
    if size == 0 {
        panic(errs.NilPointError("eureka url is empty"))
    }
    index := int(ct) % size
    url := c.eurekaUrl[index]
    //if strings.LastIndex(url,"/")>-1{
    url = strings.TrimSuffix(url, "/")
    //}
    return url
}

func (d *EurekaDiscovery) Health() (bool, string) {
    ok, desc := true, "ok"
    i := 0
    for _, u := range d.eurekaUrl {

        url, err := url.Parse(u)
        if err != nil {
            i++
            ok, desc = false, err.Error()
            continue
        }
        healthUrl := url.Scheme + "://" + url.Host + "/health"
        res, err := http.Get(healthUrl)
        if err != nil {
            i++
            ok, desc = false, err.Error()
            continue
        }
        if res == nil || res.StatusCode != http.StatusOK {
            i++
            ok, desc = true, res.Status
            continue
        }

    }

    return ok, desc

}
