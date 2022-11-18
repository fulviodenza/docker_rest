package utils

import "net/url"

// ParamsLogs are the query params of the `docker log {id}` command
var ParamsLogs map[string]string = map[string]string{
	"stdout":     "1",
	"stderr":     "1",
	"timestamps": "1",
	"details":    "1",
	"follow":     "1",
}

// ParamsLogs are the query params of the `docker container ls -a` command
var ParamList map[string]string = map[string]string{
	"all": "1",
}

// ParamsLogs are the query params of the `docker pull {image}:{tag}` command
var ParamPull map[string]string = map[string]string{
	"tag": "latest",
}

// AddQueryParams add query params to a new
// url.Values{} query and returns it
func AddQueryParams(params map[string]string) url.Values {
	query := url.Values{}
	for k, v := range params {
		query.Set(k, v)
	}
	return query
}
