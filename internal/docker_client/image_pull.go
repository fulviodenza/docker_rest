package docker_client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fulviodenza/docker_rest/internal/utils"
)

// Pull method enables to pull docker images
// specified as "image:tag", for example "ubuntu:latest"
func (dc *ClientDocker) Pull(image string) error {
	return dc.pull(image)
}

// composes a pull request over the socket
func (dc *ClientDocker) pull(image string) error {

	refs := strings.Split(image, ":")

	var q = struct {
		AttachStdout bool
		AttachStderr bool
	}{
		AttachStderr: true,
		AttachStdout: true,
	}

	params := map[string]string{
		"fromImage": refs[0],
		"tag":       refs[1],
	}
	req, err := dc.buildRequest("POST", "/v1.41/images/create", utils.AddQueryParams(params), q)
	if err != nil {
		return err
	}

	resp, err := dc.doRequest(context.Background(), req)
	if err != nil {
		return err
	}

	io.Copy(os.Stdout, resp.body)

	buf := new(strings.Builder)
	_, err = io.Copy(buf, resp.body)
	if err != nil {
		return err
	}

	fmt.Println(buf.String())
	if resp.statusCode < 200 || resp.statusCode > 299 {
		return errors.New(buf.String())
	}
	return err
}
