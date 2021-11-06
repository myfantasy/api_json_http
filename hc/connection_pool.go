package hc

import (
	"encoding/json"
	"io/ioutil"
	"sync"
	"time"

	"github.com/myfantasy/api_json"
	"github.com/myfantasy/compress"
)

type ConnectionSettings struct {
	AuthType        string          `json:"auth_type"`
	UserNameRequest string          `json:"user_name_request"`
	SecretInfo      json.RawMessage `json:"secret_info"`

	Server              string        `json:"server"`
	IgnoreSSLValidation bool          `json:"ignore_ssl_validation"`
	QueryWait           time.Duration `json:"query_wait"`
	MaxConnsPerHost     int           `json:"max_conns_per_host"`
	MaxIdleConnDuration time.Duration `json:"max_idle_conn_duration"`
}

type ConnectionPool struct {
	ConnSettings map[string]ConnectionSettings `json:"connections"`

	Connections map[string]*api_json.ApiProvider `json:"-"`

	mx sync.RWMutex
}

func (cp *ConnectionPool) InitConnections(
	compressor *compress.Generator,
	getCompression api_json.GetCompressionFunc,
) {
	cp.mx.Lock()
	defer cp.mx.Unlock()
	ap := make(map[string]*api_json.ApiProvider, len(cp.ConnSettings))

	for k, v := range cp.ConnSettings {
		conn := CreateJsonProvider(
			compressor, getCompression,
			v.AuthType, v.UserNameRequest, v.SecretInfo,
			v.Server, v.IgnoreSSLValidation,
			v.QueryWait, v.MaxConnsPerHost, v.MaxIdleConnDuration,
		)
		ap[k] = conn
	}

	cp.Connections = ap
}

func ConnectionPoolFromJson(j json.RawMessage) (cp *ConnectionPool, err error) {
	cp = &ConnectionPool{}
	err = json.Unmarshal(j, cp)
	if err != nil {
		return nil, err
	}
	return cp, nil
}

func ConnectionPoolFromFile(filename string) (cp *ConnectionPool, err error) {
	j, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return ConnectionPoolFromJson(j)
}

func (cp *ConnectionPool) Get(name string) (ap *api_json.ApiProvider, ok bool) {
	cp.mx.RLock()
	defer cp.mx.RUnlock()

	if len(cp.Connections) == 0 {
		return nil, false
	}

	ap, ok = cp.Connections[name]

	return ap, ok
}
