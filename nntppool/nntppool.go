package nntppool

import (
	"github.com/oremj/go-nntp/nntp"
	"sync"
)

type Pool struct {
	server string
	user string
	password string
	clients int
	available chan *nntp.Client
	lk sync.Mutex
}

func NewPool(server string, user string, pass string, maxClients int) *Pool {
	return &Pool{
		server: server,
		user: user,
		password: pass,
		clients: 0,
		available: make(chan *nntp.Client, maxClients),
	}
}

func (pool *Pool) grabClientLock() bool {
	pool.lk.Lock()
	defer pool.lk.Unlock()
	if pool.clients < cap(pool.available) {
		pool.clients++
		return true
	}
	return false
}

func (pool *Pool) makeClient() (client *nntp.Client) {
	if pool.grabClientLock() {
		client, _ = nntp.DialAuth(pool.server, pool.user, pool.password)
		return client
	}
	return nil
}

func (pool *Pool) GetClient() *nntp.Client {
	select {
		case conn := <-pool.available:
			return conn
		default:
			conn := pool.makeClient()
			if conn != nil {
				return conn
			}
			return <-pool.available
	}
}

func (pool *Pool) FreeClient(client *nntp.Client) {
	pool.available <-client
}
