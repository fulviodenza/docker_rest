package docker_client

import (
	"context"
	"net/url"
)

func (dc *ClientDocker) Start(id, image string) error {
	return dc.start(id, image)
}

func (dc *ClientDocker) start(image, id string) error {

	// Exec
	query := url.Values{}
	httpReq, err := dc.buildRequest("POST", "/v1.41/containers/"+id+"/start", query, struct{}{})
	if err != nil {
		return err
	}

	response, err := dc.doRequest(context.Background(), httpReq)
	if err != nil {
		return err
	}
	defer response.body.Close()

	return nil
}
