package docker_client

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/docker/docker/client"
	"github.com/docker/go-connections/sockets"
)

type serverResponse struct {
	body       io.ReadCloser
	header     http.Header
	statusCode int
	reqURL     *url.URL
}

type Client interface {
	Pull(image string) error
	Run(cmd string) error
}

type ClientDocker struct {
	// scheme sets the scheme for the client
	Scheme string
	// host holds the server address to connect to
	Host string
	// proto holds the client protocol i.e. unix.
	Proto string
	// addr holds the client address.
	Addr string
	// basePath holds the path to prepend to the requests.
	BasePath string
	// client used to send and receive http requests.
	Client *http.Client
	// custom http headers configured by users.
	CustomHTTPHeaders map[string]string
}

var _ Client = (*ClientDocker)(nil)

func NewDockerClient() (*ClientDocker, error) {

	client, err := defaultHTTPClient(DefaultDockerHost)
	if err != nil {
		return nil, err
	}

	return &ClientDocker{
		Host:   DefaultDockerHost,
		Proto:  defaultProto,
		Client: client,
	}, nil
}

func defaultHTTPClient(host string) (*http.Client, error) {
	url, err := client.ParseHostURL("http://" + host)
	if err != nil {
		return nil, err
	}
	transport := new(http.Transport)
	sockets.ConfigureTransport(transport, url.Scheme, url.Host)
	return &http.Client{
		Transport: transport,
	}, nil
}

func (dc *ClientDocker) Pull(image string) error {
	return dc.pull(image)
}

// composes a pull request over the socket
func (dc *ClientDocker) pull(image string) error {

	refs := strings.Split(image, ":")

	query := url.Values{}
	query.Set("fromImage", refs[0])
	query.Set("tag", refs[1])

	req, err := dc.buildRequest("POST", "/v1.41/images/create", query, nil)
	if err != nil {
		return nil
	}

	resp, err := dc.doRequest(context.TODO(), req)
	if err != nil {
		return nil
	}

	io.Copy(os.Stdout, resp.body)

	return nil
}

func (dc *ClientDocker) Run(cmd string) error {
	return dc.run()
}

func (dc *ClientDocker) run() error {
	// TODO: Implement same logic but for run
	// container command (possibly also with commands)
	return nil
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

func (dc *ClientDocker) buildRequest(method, path string, query url.Values, body io.Reader) (*http.Request, error) {

	url := url.URL{
		Host: dc.Host, Path: path, RawQuery: query.Encode(), Scheme: "http",
	}
	// TODO: Fix url string, currently returns
	// "http://%2Fvar%2Frun%2Fdocker.sock/v1.41/images/create?fromImage=ubuntu&tag=latest”
	// should return realistic url of the docker sock
	urlStr := url.String()
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}

	req.URL.Host = dc.Host
	return req, nil
}
