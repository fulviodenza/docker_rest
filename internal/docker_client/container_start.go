package docker_client

import (
	"context"
	"net/url"
)

// Start Method starts the container with the
// given id and image. Returns an error if the path to
// the container is not correct.
func (dc *ClientDocker) Start(id, image string) error {
	return dc.start(id, image)
}

func (dc *ClientDocker) start(image, id string) error {

	httpReq, err := dc.buildRequest("POST", "/v1.41/containers/"+id+"/start", url.Values{}, struct{}{})
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
