package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/myfantasy/api_json"
	"github.com/myfantasy/api_json_http/hc"
	"github.com/myfantasy/api_json_http/hs"
	"github.com/myfantasy/authentication/sat"
	"github.com/myfantasy/authorization/saz"
	"github.com/myfantasy/compress"
	"github.com/myfantasy/storage"
	"github.com/valyala/fasthttp"

	log "github.com/sirupsen/logrus"
)

var fListenAddress = flag.String("l", ":7499",
	"Listen address and port for example :8080 localhost:7498 etc")

var fTlsKey = flag.String("tls_key", "",
	"tls key; example `app/key.pem`")

var fTlsCert = flag.String("tls_cert", "",
	"tls certificate; example `app/cert.pem`")

var fSettingsFile = flag.String("s", "service.settings.json",
	"Settings file")

var compressor *compress.Generator
var fileStorage *storage.Generator
var autht *sat.SimpleAuthenticationChecker
var authz *saz.SimplePermissionChecker

var api *api_json.Api

var apiHandler func(ctx *fasthttp.RequestCtx)

func main() {
	flag.Parse()

	tf := new(log.TextFormatter)
	tf.FullTimestamp = true
	log.SetFormatter(tf)

	loadSettings()

	llevel, er0 := log.ParseLevel(settings.LogLevel)
	if er0 != nil {
		log.Fatal(er0)
	}
	log.SetLevel(llevel)

	// Init API
	api = &api_json.Api{}
	api.Compressor = compressor
	api.GetCompression = api_json.ZipCompressFunc
	api.PermissionChecker = authz
	api.AddAuthenticationChecker(autht)
	api.AddApi(&api_json.ServiceApi{})

	apiHandler = hs.FastHTTPHandler(api)

	// init FastHttp
	fhApi := &fasthttp.Server{
		Handler: fastHTTPHandler,
	}

	serverErrors := make(chan error, 1)
	go func() {
		if *fTlsKey != "" {
			log.Infof("Listen and serve TLS %v", *fListenAddress)
			serverErrors <- fhApi.ListenAndServeTLS(*fListenAddress,
				*fTlsCert, *fTlsKey)
		} else {
			log.Infof("Listen and serve %v", *fListenAddress)
			serverErrors <- fhApi.ListenAndServe(*fListenAddress)
		}
	}()

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Fatalf("Can`t start server; %v", err)

	case <-osSignals:
		log.Infof("Start shutdown...")
		go func() {
			if err := fhApi.Shutdown(); err != nil {
				log.Infof("Graceful shutdown did not complete in 5s : %v", err)
			}
		}()
	}

	log.Infof("Complete shutdown")

}

func fastHTTPHandler(ctx *fasthttp.RequestCtx) {
	path := string(ctx.Request.URI().Path())

	log.Tracef("enter call %v", path)

	if path == hc.ClusterMethodPath {
		apiHandler(ctx)
		return
	}

	if path == "/ping" {
		ping(ctx)
		return
	}

	notFound(ctx)
}

func notFound(ctx *fasthttp.RequestCtx) {
	ctx.Response.SetStatusCode(404)
	log.Tracef("notFound")
}

func ping(ctx *fasthttp.RequestCtx) {
	ctx.Response.SetStatusCode(200)
	log.Tracef("ping")
}
