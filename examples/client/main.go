package main

import (
	"context"
	"time"

	"github.com/myfantasy/api_json"
	"github.com/myfantasy/api_json_http/hc"
	"github.com/myfantasy/authentication/sat"
	"github.com/myfantasy/compress"

	log "github.com/sirupsen/logrus"
)

func main() {
	tf := new(log.TextFormatter)
	tf.FullTimestamp = true
	log.SetFormatter(tf)

	compressor := compress.GeneratorCreate(7)
	atReq := &sat.Request{Pwd: "123"}
	jp := hc.CreateJsonProvider(
		compressor, api_json.ZipCompressFunc, atReq.Type(),
		"admin", atReq.ToSecretInfo(),
		"http://localhost:7499",
		true,
		time.Second*5, 6, time.Second*20,
	)

	err := api_json.Ping(context.Background(), jp)

	if err != nil {
		log.Error(err)
	}
}
