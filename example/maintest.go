package main

import (
    "encoding/xml"
    "fmt"
    "io/ioutil"
    "net/http"

    "github.com/tietang/go-eureka-client/eureka"
)

func main() {
    eurekaUrl := "http://dev.discovery.shishike.com/eureka"
    //	timeout := 10 * time.Second
    //	client := &http.Client{
    //		//		CheckRedirect: redirectPolicyFunc,
    //		Timeout: timeout,
    //	}
    //	req, err := http.NewRequest("GET", eurekaUrl+"/apps", nil)
    //	//	req.Header.Add("Accept", "application/json")
    //	res, err := client.Do(req)
    res, err := http.Get(eurekaUrl + "/apps")
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(res.StatusCode)
    respBody, err := ioutil.ReadAll(res.Body)
    if err != nil {
        fmt.Println(err)
        return
    }
    if res.StatusCode != http.StatusOK {
        fmt.Println(res.StatusCode)
        return
    }
    var applications *eureka.Applications = new(eureka.Applications)
    fmt.Println(string(respBody))
    err = xml.Unmarshal(respBody, applications)
    fmt.Println(err, applications)

}
