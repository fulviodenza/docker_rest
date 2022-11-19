package docker_client

import (
	"context"
	"encoding/json"
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
		Cmd:          cmd,   // "Cmd\":[\"cat\",\"/proc/loadavg\"]
		Image:        image, // "Image\":\"ubuntu\"
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

	var respP struct {
		Id string `json:"Id"`
	}
	json.NewDecoder(resp.body).Decode(&respP)

	return respP.Id, nil
}
