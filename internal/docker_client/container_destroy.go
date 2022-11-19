package docker_client

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/fulviodenza/docker_rest/internal/utils"
)

func (dc *ClientDocker) Destroy(ctx context.Context, id string) {
	dc.destroy(ctx, id)
}

func (dc *ClientDocker) destroy(ctx context.Context, id string) {
	req, err := dc.buildRequest("DELETE", "/v1.41/containers/"+id, utils.AddQueryParams(utils.ParamList), struct{}{})
	if err != nil {
		return
	}

	resp, err := dc.doRequest(ctx, req)
	if err != nil {
		return
	}

	buf := new(strings.Builder)
	if _, err := io.Copy(buf, resp.body); err != nil {
		return
	}

	fmt.Println(buf.String())
	if resp.statusCode < 200 || resp.statusCode > 299 {
		return
	}
}
