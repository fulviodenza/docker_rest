package docker_client

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
)

type serverResponse struct {
	body       io.ReadCloser
	header     http.Header
	statusCode int
	reqURL     *url.URL
}

type Client interface {
	Pull(image string) error
	Start(id, image string) error
	Logs(id string) error
	List(ctx context.Context) (Containers, error)
	// TODO: Add Create when will fully be implemented
}

type ClientDocker struct {
	// scheme sets the scheme for the client, i.e. http
	Scheme string
	// host holds the server address to connect to, i.e. docker.sock
	Host string
	// proto holds the client protocol i.e. unix.
	Proto string
	// basePath holds the path to prepend to the requests.
	BasePath string
	// client used to send and receive http requests.
	Client *http.Client
}

var _ Client = (*ClientDocker)(nil)

// NewDockerClient return a new wrapper of the
// client connected to the docker unix socket
// this is created using `DialContext` function
// with `"unix"` and `"/var/run/docker.sock"`
// as parameters.
func NewDockerClient() *ClientDocker {

	client := defaultHTTPClient()

	return &ClientDocker{
		Host:   DefaultDockerHost,
		Proto:  defaultProto,
		Client: client,
	}
}

func defaultHTTPClient() *http.Client {
	d := new(net.Dialer)
	x := &http.Transport{
		DialContext: func(ctx context.Context, net, addr string) (net.Conn, error) {
			return d.DialContext(ctx, "unix", "/var/run/docker.sock")
		},
	}

	return &http.Client{
		Transport: x,
	}
}

func (dc *ClientDocker) doRequest(ctx context.Context, req *http.Request) (serverResponse, error) {

	serverResp := serverResponse{statusCode: -1, reqURL: req.URL}

	req = req.WithContext(ctx)
	resp, err := dc.Client.Do(req)
	if err != nil {
		return serverResp, err
	}
	if resp != nil {
		serverResp.statusCode = resp.StatusCode
		serverResp.body = resp.Body
		serverResp.header = resp.Header
	}

	return serverResp, nil
}

func (dc *ClientDocker) buildRequest(method, path string, query url.Values, req any) (*http.Request, error) {

	u := &url.URL{
		Scheme:   "http",
		Host:     dc.Host,
		Path:     path,
		RawQuery: query.Encode(),
	}

	var buf []byte
	bbuf := bytes.NewBuffer(buf)

	read, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	bbuf.Write(read)

	header := make(http.Header)
	header.Add("Content-Type", "application/json")
	httpReq := &http.Request{
		Method:     method,
		URL:        u,
		Header:     header,
		Body:       ioutil.NopCloser(bbuf),
		Host:       dc.Host,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
	}

	if bbuf.String() == "{}" {
		httpReq.Body = nil
	}

	return httpReq, err
}
