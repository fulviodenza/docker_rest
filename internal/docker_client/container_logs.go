package docker_client

import (
	"context"
	"io"
	"os"

	"github.com/fulviodenza/docker_rest/internal/utils"
)

// Logs shows logs for the container with the given id
func (dc *ClientDocker) Logs(id string) error {
	return dc.logs(id)
}

func (dc *ClientDocker) logs(id string) error {

	httpReq, err := dc.buildRequest("GET", "/containers/"+id+"/logs", utils.AddQueryParams(utils.ParamsLogs), struct{}{})
	if err != nil {
		return err
	}

	resp, err := dc.doRequest(context.Background(), httpReq)
	if err != nil {
		return err
	}
	defer resp.body.Close()

	io.Copy(os.Stdout, resp.body)

	return nil
}
