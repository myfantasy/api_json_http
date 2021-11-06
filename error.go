package api_json_http

import "github.com/myfantasy/mft"

// Errors codes and description
var Errors map[int]string = map[int]string{
	20500000: "api_json_http.hc.ConnectionFromJson: fail to unmarshal",
	20500100: "api_json_http.hc.Connection.DoRawQuery: request create error to server: `%v`",
	20500101: "api_json_http.hc.Connection.DoRawQuery: request send error to server: `%v`",
	20500102: "api_json_http.hc.Connection.DoRawQuery: read body error from request to server: `%v`",
	20500200: "api_json_http.hc.HTTPCallProvider.CallFunction: send request fail",
	20500201: "api_json_http.hc.HTTPCallProvider.CallFunction: responce code is not 200 responce code is: `%v` body: `%v`",
}

func init() {
	mft.AddErrorsCodes(Errors)
}
