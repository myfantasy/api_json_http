package hc

import (
	"encoding/json"
	"time"

	"github.com/myfantasy/api_json"
	"github.com/myfantasy/compress"
)

func CreateJsonProvider(compressor *compress.Generator,
	getCompression api_json.GetCompressionFunc,
	authType string,
	userNameRequest string,
	secretInfo json.RawMessage,

	server string,
	ignoreSSLValidation bool,
	queryWait time.Duration,
	maxConnsPerHost int,
	maxIdleConnDuration time.Duration,
) *api_json.ApiProvider {

	h := &HTTPCallProvider{
		Connection: (&Connection{
			Server:              server,
			IgnoreSSLValidation: ignoreSSLValidation,
			QueryWait:           queryWait,
			MaxConnsPerHost:     maxConnsPerHost,
			MaxIdleConnDuration: maxIdleConnDuration,
		}).Init(),
	}

	p := &api_json.ApiProvider{
		Compressor:      compressor,
		GetCompression:  getCompression,
		AuthType:        authType,
		UserNameRequest: userNameRequest,
		SecretInfo:      secretInfo,

		CallFunc: h.CallFunction,
	}

	return p
}
