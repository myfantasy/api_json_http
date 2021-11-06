package hc

import (
	"context"
	"time"

	"github.com/myfantasy/compress"
	"github.com/myfantasy/mft"
)

const (
	CompressTypeHeader = "Compress-Type"
)

var ClusterMethodPath = "/cluster"

type HTTPCallProvider struct {
	Connection *Connection
}

func (hp *HTTPCallProvider) CallFunction(ctx context.Context, compType compress.CompressionType,
	bodyRequest []byte, waitDuration time.Duration,
) (outCompType compress.CompressionType, bodyResponce []byte, err *mft.Error) {
	bodyIn, headersOut, statusCode, err := hp.Connection.DoRawQuery(waitDuration, ClusterMethodPath,
		map[string]string{CompressTypeHeader: string(compType)}, bodyRequest)

	if err != nil {
		return outCompType, nil, mft.GenerateErrorE(20500200, err)
	}

	if statusCode != 200 {
		return outCompType, nil, mft.GenerateErrorE(20500201, err, statusCode, string(bodyIn))
	}

	ct, ok := headersOut[CompressTypeHeader]
	if !ok {
		ct = string(compress.NoCompression)
	}

	return compress.CompressionType(ct), bodyIn, nil
}
