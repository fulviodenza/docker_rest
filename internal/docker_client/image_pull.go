package docker_client

import (
	"context"
	"io"
	"net/url"
	"os"
	"strings"
)

// Pull method enables to pull docker images
// specified as "image:tag", for example "ubuntu:latest"
func (dc *ClientDocker) Pull(image string) error {
	return dc.pull(image)
}

// composes a pull request over the socket
func (dc *ClientDocker) pull(image string) error {

	refs := strings.Split(image, ":")
	query := url.Values{}
	query.Set("fromImage", refs[0])
	query.Set("tag", refs[1])

	var q = struct {
		AttachStdout bool
		AttachStderr bool
	}{
		AttachStderr: true,
		AttachStdout: true,
	}

	// TODO: The "/v1.41" should be replaced to be dynamic
	req, err := dc.buildRequest("POST", "/v1.41/images/create", query, q)
	if err != nil {
		return nil
	}

	resp, err := dc.doRequest(context.Background(), req)
	if err != nil {
		return nil
	}

	io.Copy(os.Stdout, resp.body)

	return err
}
