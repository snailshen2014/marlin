package codisclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/op/go-logging"
	"github.com/qiniu/x/log.v7"
	"github.com/samuel/go-zookeeper/zk"
)

type proxyInfo struct {
	ProtoType string `json:"proto_type"`
	ProxyAddr string `json:"proxy_addr"`
}

//Pool connection pool
type Pool struct {
	nextIdx int
	//added by shenyanjun for redis.Pool
	// pools     []redis.Conn
	pools     []*redis.Pool
	zk        zk.Conn
	ZkServers []string
	ZkTimeout time.Duration
	ZkDir     string
	Dial      func(network, address string) (redis.Conn, error)
}

var (
	mu                  = &sync.Mutex{}
	Log *logging.Logger = logging.MustGetLogger("Codis")
)

func (pool *Pool) initFromZk() error {
	pool.initZk()
	// this.pools = []redis.Conn{}
	pool.pools = []*redis.Pool{}
	children, _, err := pool.zk.Children(pool.ZkDir)

	if err != nil {
		//Fatal 会使程序退出
		// Log.Fatal(err)ß
		log.Error(err)
		return err
	}

	for _, child := range children {
		data, _, err := pool.zk.Get(pool.ZkDir + "/" + child)
		if err != nil {
			continue
		}

		var p proxyInfo

		json.Unmarshal(data, &p)
		// conn, err := this.Dial(p.Proto_type, p.Proxy_addr)
		// if err != nil {
		// 	Log.Errorf("Create redis connection failed: %s", err.Error())
		// 	continue
		// }
		one := initRedisPool(p.ProxyAddr)
		fmt.Printf("pool get a pool to:%s\n", p.ProxyAddr)
		pool.pools = append(pool.pools, one)
	}

	go pool.watch(pool.ZkDir)
	return nil
}

func (pool *Pool) watch(node string) {
	for {
		_, _, ch, err := pool.zk.ChildrenW(node)
		if err != nil {
			Log.Error(err)
			return
		}
		evt := <-ch

		if evt.Type == zk.EventSession {
			if evt.State == zk.StateConnecting {
				continue
			}
			if evt.State == zk.StateExpired {
				pool.zk.Close()
				Log.Info("Zookeeper session expired, reconnecting...")
				pool.initZk()
			}
		}
		if evt.State == zk.StateConnected {
			switch evt.Type {
			case
				zk.EventNodeCreated,
				zk.EventNodeDeleted,
				zk.EventNodeChildrenChanged,
				zk.EventNodeDataChanged:
				err := pool.initFromZk()
				if err != nil {
					log.Error("initFromZk error:" + err.Error())
				}
				return
			}
			continue
		}
	}
}

func (pool *Pool) initZk() {
	zkConn, _, err := zk.Connect(pool.ZkServers, pool.ZkTimeout)
	if err != nil {
		Log.Fatalf("Failed to connect to zookeeper: %+v", err)
	}
	pool.zk = *zkConn
}

//Get get a connection from one pool
func (pool *Pool) Get() (redis.Conn, error) {
	mu.Lock()
	if len(pool.pools) == 0 {
		err := pool.initFromZk()
		if err != nil {
			mu.Unlock()
			log.Error("initFromZk error:" + err.Error())
			return nil, err
		}
	}
	pool.nextIdx++
	if pool.nextIdx >= len(pool.pools) {
		pool.nextIdx = 0
	}
	if len(pool.pools) == 0 {
		mu.Unlock()
		err := errors.New("Proxy list empty")
		Log.Error(err)
		return errorConnection{err: err}, nil
	}
	one := pool.pools[pool.nextIdx]
	mu.Unlock()
	// return c
	return one.Get(), nil

}

//Close close pool
func (pool *Pool) Close() {
	pool.zk.Close()
	// this.pools = []redis.Conn{}
	pool.pools = []*redis.Pool{}
}

type errorConnection struct{ err error }

func (ec errorConnection) Do(string, ...interface{}) (interface{}, error) { return nil, ec.err }
func (ec errorConnection) Send(string, ...interface{}) error              { return ec.err }
func (ec errorConnection) Err() error                                     { return ec.err }
func (ec errorConnection) Close() error                                   { return ec.err }
func (ec errorConnection) Flush() error                                   { return ec.err }
func (ec errorConnection) Receive() (interface{}, error)                  { return nil, ec.err }

//InitRedisPool ,codis.Pool没有实现真正的redis线程池，所以用redi自身的线程池替代
func initRedisPool(host string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     16,
		MaxActive:   16,
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", host)
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}
}
