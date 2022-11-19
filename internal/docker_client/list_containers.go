package docker_client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/fulviodenza/docker_rest/internal/utils"
)

// List method lists all (`--all`) containers on the machine
func (dc *ClientDocker) List(ctx context.Context) (Containers, error) {
	return dc.list(ctx)
}

func (dc *ClientDocker) list(ctx context.Context) (Containers, error) {

	req, err := dc.buildRequest("GET", "/v1.41/containers/json", utils.AddQueryParams(utils.ParamList), struct{}{})
	if err != nil {
		return nil, err
	}

	resp, err := dc.doRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	containers := Containers{}

	if resp.statusCode < 200 || resp.statusCode > 299 {
		return Containers{}, errors.New("Got status code: " + fmt.Sprint(resp.statusCode))
	}

	return containers, json.NewDecoder(resp.body).Decode(&containers)
}
