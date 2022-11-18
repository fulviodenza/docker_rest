package docker_client

import (
	"context"
	"io"
	"net/url"
	"os"
)

func (dc *ClientDocker) Logs(id string) error {
	return dc.logs(id)
}

func (dc *ClientDocker) logs(id string) error {
	// TODO: utility function to set query parameters
	query := url.Values{}
	query.Set("stdout", "1")
	query.Set("stderr", "1")
	query.Set("timestamps", "1")
	query.Set("details", "1")
	query.Set("follow", "1")

	httpReq, err := dc.buildRequest("GET", "/containers/"+id+"/logs", query, struct{}{})
	if err != nil {
		return err
	}

	response, err := dc.doRequest(context.Background(), httpReq)
	if err != nil {
		return err
	}
	io.Copy(os.Stdout, response.body)

	return nil

}
