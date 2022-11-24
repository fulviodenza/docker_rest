package docker_client

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

type serverResponse struct {
	body       io.ReadCloser
	header     http.Header
	statusCode int
	reqURL     *url.URL
}

// Client is the interface to execute main operation
// towards the docker daemon.
type Client interface {
	List(ctx context.Context) (Containers, error)
	Pull(image string) error
	Create(image string, cmd []string) (string, error)
	Start(id string) error
	Logs(id string) error
	Destroy(ctx context.Context, id string)
}

// ClientDocker is the structure which wraps the http client
// and contains specification for executing the request, for example
// the Scheme, the Host and the Proto.
type ClientDocker struct {
	// scheme sets the scheme for the client, i.e. http
	Scheme string
	// host holds the server address to connect to, i.e. docker.sock
	Host string
	// proto holds the client protocol i.e. unix.
	Proto string
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
		Proto:  DefaultProto,
		Scheme: DefaultScheme,
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
		Timeout:   DefaultTimeout * time.Second,
	}
}

// buildRequest is the method to build the request starting
// from a method, a path, a query and a struct{} containing
// the body of the request. This function should be coupled
// with doRequest which wraps the Do method of the http.Client
func (dc *ClientDocker) buildRequest(method, path string, query url.Values, req any) (*http.Request, error) {

	u := &url.URL{
		Scheme:   dc.Scheme,
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
		Body:       io.NopCloser(bbuf),
		Host:       dc.Host,
		Proto:      dc.Proto,
		ProtoMajor: 1,
		ProtoMinor: 1,
	}

	if bbuf.String() == "{}" {
		httpReq.Body = nil
	}

	return httpReq, err
}

// doRequest wraps the Do method of the Client wrapped
// inside the ClientDocker structure. It do the request
// it gets from buildRequest and returns a serverResponse
// with a StatuCode, a Body, and a header.
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
