package docker_client

import (
	"context"
	"errors"
	"io"
	"net/url"
	"strings"
)

// Start Method starts the container with the
// given id and image. Returns an error if the path to
// the container is not correct.
func (dc *ClientDocker) Start(id string) error {
	return dc.start(context.Background(), id)
}

func (dc *ClientDocker) start(ctx context.Context, id string) error {

	query := url.Values{}
	httpReq, err := dc.buildRequest("POST", "/v1.41/containers/"+id+"/start", query, struct{}{})
	if err != nil {
		return err
	}

	resp, err := dc.doRequest(context.Background(), httpReq)
	if err != nil {
		return err
	}
	buf := new(strings.Builder)
	_, err = io.Copy(buf, resp.body)
	if err != nil {
		return err
	}

	if resp.statusCode < 200 || resp.statusCode > 299 {
		return errors.New(buf.String())
	}

	return nil
}
