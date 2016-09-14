package zuul

//import (
//	"net/http"
//	"time"

//	"github.com/jolestar/go-commons-pool"
//)

//type HttpClient struct {
//	pool    *pool.ObjectPool
//	timeout time.Duration
//}

//func (h *HttpClient) Init(config *pool.ObjectPoolConfig) {
//	if config == nil {
//		config = pool.NewDefaultPoolConfig()
//	}
//	abandonedConfig := pool.NewDefaultAbandonedConfig()
//	pooledObjectFactory := pool.NewPooledObjectFactory(h.create, destroy, validate, activate, passivate)
//	h.pool = pool.NewObjectPoolWithAbandonedConfig(pooledObjectFactory, config, abandonedConfig)

//}

//func (h *HttpClient) create() (interface{}, error) {
//	client := &http.Client{
//		Timeout: h.timeout,
//	}
//	return client, nil

//}
//func destroy(object *PooledObject) error {

//}
//func validate(object *PooledObject) bool {

//}
//func activate(object *PooledObject) error {

//}
//func passivate(object *PooledObject) error {

//}
