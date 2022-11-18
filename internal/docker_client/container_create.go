package docker_client

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/url"
)

// Create, given an image and a cmd, creates
// a new container with the given Cmd, the given Image
// and sets AttachStdout and AttachStderr true
func (dc *ClientDocker) Create(image string, cmd []string) (string, error) {
	ctx := context.Background()
	return dc.create(ctx, image, cmd)
}

func (dc *ClientDocker) create(ctx context.Context, image string, cmd []string) (string, error) {

	var req = struct {
		Cmd          []string
		Image        string
		AttachStdout bool
		AttachStderr bool
	}{
		Cmd:          cmd, // "Cmd\":[\"cat\",\"/proc/loadavg\"],\"Image\":\"ubuntu\"
		Image:        image,
		AttachStderr: true,
		AttachStdout: true,
	}

	httpReq, err := dc.buildRequest("POST", "/v1.41/containers/create", url.Values{}, req)
	if err != nil {
		return "", err
	}

	resp, err := dc.doRequest(ctx, httpReq)
	if err != nil {
		return "", err
	}
	defer resp.body.Close()

	// Parse the response inside server structure
	bodyBytes, err := ioutil.ReadAll(resp.body)
	if err != nil {
		return "", err
	}

	var respP struct {
		Id string `json:"Id"`
	}
	if err := json.Unmarshal(bodyBytes, &respP); err != nil {
		return "", err
	}

	return respP.Id, nil
}
