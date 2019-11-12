package httpclient

import (
    "bytes"
    "fmt"
    log "github.com/sirupsen/logrus"
    "github.com/tietang/props/kvs"
    "io/ioutil"
    "net"
    "net/http"
    "sync"
    "time"
)

//
//http.Client.connect.timeout=3000

const (
    DEFAULT_CONNECT_TIMEOUT         = 30 * time.Second
    DEFAULT_MAX_IDLE_CONNS          = 512
    DEFAULT_MAX_IDLE_CONNS_PER_HOST = 10
    //
    DEFAULT_IDLE_CONN_TIMEOUT       = 30 * time.Second
    DEFAULT_RESPONSE_HEADER_TIMEOUT = 10 * time.Second
    DEFAULT_KEEP_ALIVE              = 30 * time.Second
    DEFAULT_EXPECT_CONTINUE_TIMEOUT = 10 * time.Second
)

const (
    HTTP_CONCONNECT_TIMEOUT  = "http.client.connect.timeout"
    HTTP_CONCONNECT_PREFIX   = "%shttp.client."
    HTTP_CONNECTT_IMEOUT     = HTTP_CONCONNECT_PREFIX + "connectTimeout"
    HTTP_ReadTimeout         = HTTP_CONCONNECT_PREFIX + "readTimeout"
    HTTP_WriteTimeout        = HTTP_CONCONNECT_PREFIX + "writeTimeout"
    HTTP_MaxConnDuration     = HTTP_CONCONNECT_PREFIX + "maxConnDuration"
    HTTP_MaxIdleConnDuration = HTTP_CONCONNECT_PREFIX + "maxIdleConnDuration"
    HTTP_MaxResponseBodySize = HTTP_CONCONNECT_PREFIX + "maxResponseBodySize"
    HTTP_WriteBufferSize     = HTTP_CONCONNECT_PREFIX + "writeBufferSize"
    HTTP_ReadBufferSize      = HTTP_CONCONNECT_PREFIX + "readBufferSize"
    HTTP_MaxConns            = HTTP_CONCONNECT_PREFIX + "maxConns"
    //
    HTTP_KEEPALIVE               = HTTP_CONCONNECT_PREFIX + "keepAlive"
    HTTP_MAX_IDLE_CONNS          = HTTP_CONCONNECT_PREFIX + "maxIdleConns"
    HTTP_MAX_IDLE_CONNS_PER_HOST = HTTP_CONCONNECT_PREFIX + "maxIdleConnsPerHost"
    HTTP_IDLE_CONN_TIMEOUT       = HTTP_CONCONNECT_PREFIX + "idleConnTimeout"
    HTTP_RESPONSE_HEADER_TIMEOUT = HTTP_CONCONNECT_PREFIX + "responseHeaderTimeout"
    HTTP_EXPECT_CONTINUE_TIMEOUT = HTTP_CONCONNECT_PREFIX + "expectContinueTimeout"
)

type HttpClient struct {
    Name         string
    Client       *http.Client
    ClientConfig *HttpClientConfig
}

//Millisecond
type HttpClientConfig struct {
    ConnectTimeout time.Duration
    KeepAlive      time.Duration
    //
    MaxIdleConns          int
    MaxIdleConnsPerHost   int
    IdleConnTimeout       time.Duration
    ResponseHeaderTimeout time.Duration
    ExpectContinueTimeout time.Duration
}

func copyHttpClientConfig(hcc *HttpClientConfig) *HttpClientConfig {
    return &HttpClientConfig{
        ConnectTimeout:        hcc.ConnectTimeout,
        KeepAlive:             hcc.KeepAlive,
        MaxIdleConns:          hcc.MaxIdleConns,
        MaxIdleConnsPerHost:   hcc.MaxIdleConnsPerHost,
        IdleConnTimeout:       hcc.IdleConnTimeout,
        ResponseHeaderTimeout: hcc.ResponseHeaderTimeout,
        ExpectContinueTimeout: hcc.ExpectContinueTimeout,
    }

}

type HttpClients struct {
    defaultClientConfig *HttpClientConfig
    conf                kvs.ConfigSource
    clients             map[string]*HttpClient

    lock sync.Mutex
}

func NewHttpClients(conf kvs.ConfigSource) *HttpClients {

    hcs := &HttpClients{
        clients: make(map[string]*HttpClient),
        conf:    conf,
    }
    defaultClientConfig := &HttpClientConfig{}
    config(defaultClientConfig, "", conf)
    hcs.defaultClientConfig = defaultClientConfig

    return hcs
}

func config(hcc *HttpClientConfig, name string, conf kvs.ConfigSource) {
    hcc.MaxIdleConns = conf.GetIntDefault(fmt.Sprintf(HTTP_MAX_IDLE_CONNS, name), DEFAULT_MAX_IDLE_CONNS)
    hcc.MaxIdleConnsPerHost = conf.GetIntDefault(fmt.Sprintf(HTTP_MAX_IDLE_CONNS_PER_HOST, name), DEFAULT_MAX_IDLE_CONNS_PER_HOST)
    hcc.ConnectTimeout = conf.GetDurationDefault(fmt.Sprintf(HTTP_CONNECTT_IMEOUT, name), DEFAULT_CONNECT_TIMEOUT)
    hcc.ResponseHeaderTimeout = conf.GetDurationDefault(fmt.Sprintf(HTTP_RESPONSE_HEADER_TIMEOUT, name), DEFAULT_RESPONSE_HEADER_TIMEOUT)
    hcc.KeepAlive = conf.GetDurationDefault(fmt.Sprintf(HTTP_KEEPALIVE, name), DEFAULT_KEEP_ALIVE)
    hcc.IdleConnTimeout = conf.GetDurationDefault(fmt.Sprintf(HTTP_IDLE_CONN_TIMEOUT, name), DEFAULT_IDLE_CONN_TIMEOUT)
    hcc.ExpectContinueTimeout = conf.GetDurationDefault(fmt.Sprintf(HTTP_EXPECT_CONTINUE_TIMEOUT, name), DEFAULT_EXPECT_CONTINUE_TIMEOUT)
}

func (h *HttpClients) GetOrCreate(name string) *HttpClient {
    h.lock.Lock()
    defer h.lock.Unlock()
    c, ok := h.clients[name]
    if ok {
        return c
    }
    client := &HttpClient{}
    clientConfig := copyHttpClientConfig(h.defaultClientConfig)
    config(clientConfig, name, h.conf)
    client.ClientConfig = clientConfig

    hcc := clientConfig
    //fmt.Println(hcc.ConnectTimeout)
    //http client设置
    transport := &http.Transport{
        Dial: func(network, addr string) (net.Conn, error) {
            //deadline := time.Now().Add(10 * time.Second)
            c, err := net.DialTimeout(network, addr, hcc.ConnectTimeout) //设置建立连接超时
            if err != nil {
                return nil, err
            }
            //c.SetDeadline(deadline) //不建议，设置发送接收数据超时
            return c, nil
        },
        //DialContext: (&net.Dialer{
        //    Timeout:   hcc.ConnectTimeout, // 30s  time.Duration(config.ConnectTimeout) * time.Millisecond,
        //    KeepAlive: hcc.KeepAlive,      //30s,
        //    DualStack: true,
        //}).DialContext,
        //MaxIdleConns:          hcc.MaxIdleConns,          //100,最大空闲(keep-alive)连接，
        //MaxIdleConnsPerHost:   hcc.MaxIdleConnsPerHost,   //2,每host最大空闲(keep-alive)连接，=0则DefaultMaxIdleConnsPerHost=2
        //IdleConnTimeout:       hcc.IdleConnTimeout,       //90s,90 * time.Second,
        //ExpectContinueTimeout: hcc.ExpectContinueTimeout, //10s,time.Second * 2,
        //ResponseHeaderTimeout: hcc.ResponseHeaderTimeout, //
        //TLSHandshakeTimeout:   1 * time.Second,           //1s
    }
    httpClient := &http.Client{
        Transport: transport,
    }
    client.Client = httpClient
    h.clients[name] = client
    return client
}

func (h *HttpClients) Do(name string, req *http.Request) (*http.Response, error) {
    hc := h.GetOrCreate(name)

    //调用请求
    res, err := hc.Client.Do(req)

    if err != nil {
        //抛出到最外层输出
        //log.Error(err)
        return nil, err
    }
    // 如果出错就不需要close，因此defer语句放在err处理逻辑后面
    // Clients must call resp.Body.Close when finished reading resp.Body -- from golang doc
    //defer res.Body.Close()

    //处理response
    //respBody, err := ioutil.ReadAll(res.Body)
    // Reset resp.Body so it can be use again
    //res.Body = ioutil.NopCloser(bytes.NewBuffer(respBody))

    //
    //if err := res.Body.Close(); err != nil {
    //    log.Error(err)
    //}

    return res, err
}

func (h *HttpClients) DoAndRead(name string, req *http.Request) (*http.Response, []byte, error) {
    hc := h.GetOrCreate(name)

    //调用请求
    res, err := hc.Client.Do(req)

    if err != nil {
        log.Error(err)
        return nil, nil, err
    }
    // 如果出错就不需要close，因此defer语句放在err处理逻辑后面
    // Clients must call resp.Body.Close when finished reading resp.Body -- from golang doc
    defer res.Body.Close()

    //处理response
    respBody, err := ioutil.ReadAll(res.Body)
    // Reset resp.Body so it can be use again
    res.Body = ioutil.NopCloser(bytes.NewBuffer(respBody))

    //
    if err := res.Body.Close(); err != nil {
        log.Error(err)
    }

    return res, respBody, err
}
