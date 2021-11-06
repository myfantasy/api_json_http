package main

import (
	"context"
	"flag"

	"github.com/myfantasy/api_json"
	"github.com/myfantasy/api_json_http/hc"
	"github.com/myfantasy/compress"

	log "github.com/sirupsen/logrus"
)

var fSettingsFile = flag.String("cns", "connections.json",
	"Connection pool file")

func main() {
	tf := new(log.TextFormatter)
	tf.FullTimestamp = true
	log.SetFormatter(tf)

	flag.Parse()
	compressor := compress.GeneratorCreate(7)

	cp, er0 := hc.ConnectionPoolFromFile(*fSettingsFile)
	if er0 != nil {
		log.Fatal(er0)
	}

	cp.InitConnections(compressor, api_json.ZipCompressFunc)

	ap, ok := cp.Get("admin")
	if !ok {
		log.Fatal("Connection `admin` should be exists")
	}

	err := api_json.Ping(context.Background(), ap)

	if err != nil {
		log.Error(err)
	}
}
