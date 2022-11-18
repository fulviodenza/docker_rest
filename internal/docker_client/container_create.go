package docker_client

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/url"
)

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
		// "Cmd\":[\"cat\",\"/proc/loadavg\"],\"Image\":\"ubuntu\"
		Cmd:          cmd,
		Image:        image,
		AttachStderr: true,
		AttachStdout: true,
	}

	//"Cmd\":[\"cat\",\"/proc/loadavg\"],\"Image\":\"ubuntu\"
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
