package hs

import (
	"context"

	"github.com/myfantasy/api_json"
	"github.com/myfantasy/api_json_http/hc"
	"github.com/myfantasy/compress"
	"github.com/valyala/fasthttp"
)

// FastHTTPHandler - fasthttp fast http handler
func FastHTTPHandler(api *api_json.Api,
) func(ctx *fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		decompressAlg := ""
		ctx.Request.Header.VisitAll(func(key []byte, value []byte) {
			if string(key) == hc.CompressTypeHeader {
				decompressAlg = string(value)
			}
		})

		outCompType, bodyResponce := api.Do(context.Background(), compress.CompressionType(decompressAlg),
			ctx.Request.Body())

		htmlCode := 200

		ctx.Response.SetStatusCode(htmlCode)
		ctx.Response.Header.Add(hc.CompressTypeHeader, string(outCompType))
		ctx.Response.SetBody(bodyResponce)
	}
}
